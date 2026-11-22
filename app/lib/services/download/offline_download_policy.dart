import 'package:ft_music/plugins/utils/media_id.dart';

/// Fail-closed policy for saving plugin tracks for offline playback.
///
/// A content resolver must explicitly declare the `offline_download` manifest
/// capability. Direct URLs and locally imported tracks are not download sources.
class OfflineDownloadPolicy {
  OfflineDownloadPolicy._();

  static const String capability = 'offline_download';

  static bool allows({
    required String mediaId,
    required Iterable<String> capabilities,
  }) {
    final parts = tryParseMediaId(mediaId);
    if (parts == null || parts.pluginId == kLocalPluginId) return false;

    return capabilities.any(
      (value) => value.trim().toLowerCase() == capability,
    );
  }
}
