import 'dart:io';
import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';

void main() {
  test('canonical FT-music launcher mark is a valid square PNG', () async {
    final png = File('assets/icons/ft_music_mark.png');
    expect(await png.exists(), isTrue);

    final bytes = await png.readAsBytes();
    expect(bytes.take(8).toList(), <int>[
      0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
    ]);

    final header = ByteData.sublistView(Uint8List.fromList(bytes), 16, 24);
    expect(header.getUint32(0, Endian.big), 256);
    expect(header.getUint32(4, Endian.big), 256);
  });

  test('legacy duplicate icon is no longer shipped', () {
    expect(File('assets/icons/ft_music_lighter.png').existsSync(), isFalse);
  });

  test('vector mark uses the FT-music lime brand accent', () async {
    final svg = await File('assets/icons/ft_music_mark.svg').readAsString();
    expect(svg, contains('#A3E635'));
  });
}
