/// Stream resolution: selects the best stream for a download task by calling
/// the plugin manager, caching the result, and refreshing only when the cached
/// URL has expired.
use std::path::Path;
use std::sync::{Arc, Mutex};
use std::sync::atomic::Ordering;

use crate::api::downloader::types::ManagedTask;
use crate::api::downloader::utils::{
    current_unix_epoch_secs, guess_extension_from_url, normalize_quality, split_media_id,
    build_download_file_name,
};
use crate::api::plugin::commands::{ContentResolverCommand, PluginRequest, PluginResponse};
use crate::api::plugin::models::{Quality, StreamSource};
use crate::api::plugin::plugin::PluginManager;
use crate::api::plugin::types::PluginType;

/// Seconds before a stream's stated expiry at which we consider it too close to
/// use. This prevents starting a download with a URL that will expire mid-flight.
const EXPIRY_MARGIN_SECS: u64 = 60;

/// Return `true` if `stream` has a non-empty HTTPS/HTTP URL and is not expired
/// (with a safety margin of [`EXPIRY_MARGIN_SECS`]).
pub fn is_stream_usable(stream: &StreamSource) -> bool {
    let url = stream.url.trim();
    if url.is_empty() || !(url.starts_with("http://") || url.starts_with("https://")) {
        return false;
    }
    if let Some(expires_at) = stream.expires_at {
        let now = current_unix_epoch_secs();
        if expires_at <= now + EXPIRY_MARGIN_SECS {
            return false;
        }
    }
    true
}

/// A source is eligible for offline saving only when its installed content
/// resolver manifest explicitly declares this capability.
fn plugin_allows_offline_download(capabilities: &[String]) -> bool {
    capabilities
        .iter()
        .any(|capability| capability.trim().eq_ignore_ascii_case("offline_download"))
}

/// Pick the best stream from `streams` given a quality preference string.
/// Falls back to progressively wider selections before giving up.
pub fn select_best_stream(streams: &[StreamSource], preference: &str) -> Option<StreamSource> {
    let priority: &[Quality] = match normalize_quality(preference) {
        "Low" => &[Quality::Low, Quality::Medium, Quality::High, Quality::Lossless],
        "High" => &[Quality::Lossless, Quality::High, Quality::Medium, Quality::Low],
        _ => &[Quality::Medium, Quality::High, Quality::Low, Quality::Lossless],
    };

    for &quality in priority {
        if let Some(s) = streams.iter().find(|s| s.quality == quality && is_stream_usable(s)) {
            return Some(s.clone());
        }
    }
    // Last resort: any usable stream regardless of quality.
    streams.iter().find(|s| is_stream_usable(s)).cloned()
}

/// Resolve the best stream for `task_id`, using the cached stream when still
/// valid or fetching a fresh one from the plugin manager otherwise.
///
/// Also derives and stores `file_name` / `target_path` on first successful
/// resolution so the rest of the pipeline always has a non-empty target.
///
/// Returns `Err` if the task was cancelled while waiting on the plugin call.
pub async fn resolve_stream(
    task: &Arc<Mutex<ManagedTask>>,
    plugin_manager: &PluginManager,
) -> Result<StreamSource, String> {
    // 1. Return cached stream if it is still usable.
    {
        let guard = task.lock().map_err(|_| "Task mutex poisoned".to_string())?;
        if guard.cancel_requested.load(Ordering::Acquire) {
            return Err("Task was cancelled during stream resolution".to_string());
        }
        if let Some(cached) = &guard.persisted.selected_stream {
            if is_stream_usable(cached) {
                return Ok(cached.clone());
            }
        }
    }

    // 2. Need a fresh stream – read the fields we need under a short lock.
    let (track, preference) = {
        let guard = task.lock().map_err(|_| "Task mutex poisoned".to_string())?;
        (guard.persisted.track.clone(), guard.persisted.preferred_quality.clone())
    };

    // 3. Fetch from plugin manager (this is async and may take time).
    let stream = if let Some((plugin_id, local_id)) = split_media_id(&track.id) {
        let plugin_info = plugin_manager
            .get_available_plugins()
            .await
            .into_iter()
            .find(|plugin| {
                plugin.id() == plugin_id.as_str()
                    && plugin.plugin_type == PluginType::ContentResolver
            })
            .ok_or_else(|| format!("Content resolver '{plugin_id}' is not installed"))?;

        if !plugin_allows_offline_download(&plugin_info.manifest.capabilities) {
            return Err(format!(
                "Content resolver '{plugin_id}' has not explicitly authorized offline downloads. Its manifest must declare the 'offline_download' capability."
            ));
        }

        let response = plugin_manager
            .handle_plugin_request(
                &plugin_id,
                PluginRequest::ContentResolver(ContentResolverCommand::GetStreams {
                    id: local_id,
                }),
            )
            .await
            .map_err(|e| e.to_string())?;

        let streams = match response {
            PluginResponse::Streams(s) => s,
            _ => return Err("Unexpected plugin response for GetStreams".to_string()),
        };

        select_best_stream(&streams, &preference)
            .ok_or_else(|| format!("No usable streams returned for '{}'", track.title))?
    } else {
        return Err(format!(
            "Offline saving for '{}' requires a content resolver that declares the 'offline_download' capability.",
            track.title
        ));
    };

    // 4. Check for cancellation again before mutating the task (avoids a race
    //    where the user cancelled while the plugin call was in flight).
    {
        let guard = task.lock().map_err(|_| "Task mutex poisoned".to_string())?;
        if guard.cancel_requested.load(Ordering::Acquire) {
            return Err("Task was cancelled during stream resolution".to_string());
        }
    }

    // 5. Persist the resolved stream and derive file name / target path if not set.
    {
        let mut guard = task.lock().map_err(|_| "Task mutex poisoned".to_string())?;
        guard.persisted.selected_stream = Some(stream.clone());

        if guard.persisted.file_name.is_empty() || guard.persisted.target_path.is_empty() {
            let ext = guess_extension_from_url(
                Some(&stream.format),
                &stream.url,
            );
            let file_name = build_download_file_name(&guard.persisted.track, &ext);
            let target_path = Path::new(&guard.persisted.download_dir)
                .join(&file_name)
                .to_string_lossy()
                .into_owned();
            guard.persisted.file_name = file_name;
            guard.persisted.target_path = target_path;
        }
    }

    Ok(stream)
}



#[cfg(test)]
mod offline_download_policy_tests {
    use super::plugin_allows_offline_download;

    #[test]
    fn requires_explicit_offline_download_capability() {
        assert!(plugin_allows_offline_download(&[
            "charts".to_string(),
            "offline_download".to_string(),
        ]));
        assert!(!plugin_allows_offline_download(&["streaming".to_string()]));
        assert!(!plugin_allows_offline_download(&[]));
    }

    #[test]
    fn capability_matching_is_trimmed_and_case_insensitive() {
        assert!(plugin_allows_offline_download(&[
            " Offline_Download ".to_string(),
        ]));
    }
}
