# Offline music saving

FT-music keeps the existing download queue, restart recovery, cancellation, metadata tagging, and offline playback path. A download is accepted only when the installed **content resolver** explicitly opts in.

## Source authorization

A plugin's `manifest.json` must include the exact capability string:

```json
{
  "capabilities": ["offline_download"]
}
```

The app checks this declaration before queueing and again in the Rust stream resolver before requesting a stream. Missing, malformed, direct-URL, and local-file sources are denied. The capability is a declaration by the plugin publisher; users should install plugins only from publishers they trust. It is not a mechanism for bypassing a provider's terms or DRM.

## Storage and permissions

- Android and iOS use FT-music's app-scoped application-support directory for downloaded audio. This avoids requesting broad media-library or shared-storage access just to save app-managed offline files.
- Desktop uses the configured download folder (or the application documents folder when none is configured).
- Before queueing, FT-music creates the folder and performs a small write/delete check. An inaccessible folder fails with a storage-access message rather than leaving a broken task.
- Existing imported local tracks and their file permissions are unchanged.

## Queue and playback

The Rust download manager emits progress snapshots, supports pause/resume/cancel, persists task state for restart recovery, and writes audio tags and artwork on a best-effort basis. Completed files are indexed in the downloads playlist. Playback resolves a valid local download before attempting an online plugin stream, so saved tracks continue to play without network access.

Metadata/artwork tagging is best-effort: a provider may not supply a usable artwork URL or a container may not support all tags. The audio file is retained if tagging fails.
