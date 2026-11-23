import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:ft_music/core/theme/app_motion.dart';

void main() {
  group('AppMotion', () {
    testWidgets('resolves durations normally when motion is enabled',
        (tester) async {
      late Duration duration;
      await tester.pumpWidget(
        MaterialApp(
          home: Builder(
            builder: (context) {
              duration = AppMotion.resolve(context, AppMotion.standard);
              return const SizedBox.shrink();
            },
          ),
        ),
      );

      expect(duration, AppMotion.standard);
    });

    testWidgets('removes motion when accessibility disables animations',
        (tester) async {
      late Duration duration;
      await tester.pumpWidget(
        MaterialApp(
          home: MediaQuery(
            data: const MediaQueryData(disableAnimations: true),
            child: Builder(
              builder: (context) {
                duration = AppMotion.resolve(context, AppMotion.standard);
                return const SizedBox(key: ValueKey('motion-target'));
              },
            ),
          ),
        ),
      );

      expect(
        AppMotion.reduceMotion(tester.element(find.byKey(const ValueKey('motion-target')))),
        isTrue,
      );
      expect(duration, Duration.zero);
    });
  });
}
