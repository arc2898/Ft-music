import 'package:flutter_test/flutter_test.dart';
import 'package:ft_music/services/download/offline_download_policy.dart';

void main() {
  group('OfflineDownloadPolicy', () {
    test('allows a plugin track only with explicit capability', () {
      expect(
        OfflineDownloadPolicy.allows(
          mediaId: 'licensed-source::track-123',
          capabilities: const ['offline_download'],
        ),
        isTrue,
      );
    });

    test('denies plugin tracks without the capability', () {
      expect(
        OfflineDownloadPolicy.allows(
          mediaId: 'source::track-123',
          capabilities: const ['streaming'],
        ),
        isFalse,
      );
    });

    test('denies local imports and unscoped direct URLs', () {
      expect(
        OfflineDownloadPolicy.allows(
          mediaId: 'local::file-hash',
          capabilities: const ['offline_download'],
        ),
        isFalse,
      );
      expect(
        OfflineDownloadPolicy.allows(
          mediaId: 'https://example.test/audio.mp3',
          capabilities: const ['offline_download'],
        ),
        isFalse,
      );
    });

    test('requires the exact capability name', () {
      expect(
        OfflineDownloadPolicy.allows(
          mediaId: 'source::track-123',
          capabilities: const ['download'],
        ),
        isFalse,
      );
    });
  });
}
