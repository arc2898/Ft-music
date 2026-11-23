import 'package:flutter/material.dart';
import 'package:ft_music/core/theme/app_motion.dart';

/// Staggered enter transition shared by library and search result lists.
/// Uses short, capped delays and respects the platform reduced-motion setting.
class AnimatedListItem extends StatelessWidget {
  final int index;
  final Widget child;

  const AnimatedListItem({
    super.key,
    required this.index,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    final stagger = (index.clamp(0, 12) * 16);
    final duration = AppMotion.resolve(
      context,
      Duration(milliseconds: AppMotion.standard.inMilliseconds + stagger),
    );
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0, end: 1),
      duration: duration,
      curve: AppMotion.emphasized,
      builder: (context, progress, child) {
        final offset = AppMotion.reduceMotion(context) ? 0.0 : 10 * (1 - progress);
        final scale = AppMotion.reduceMotion(context) ? 1.0 : .992 + (.008 * progress);
        return Opacity(
          opacity: progress,
          child: Transform.translate(
            offset: Offset(0, offset),
            child: Transform.scale(scale: scale, child: child),
          ),
        );
      },
      child: child,
    );
  }
}
