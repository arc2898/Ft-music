import 'package:flutter/material.dart';

/// Shared motion tokens for FT-music.
///
/// Durations are deliberately short and non-oscillating. Every UI transition
/// should resolve its duration through [resolve] so Android/iOS accessibility
/// settings that disable animations are respected consistently.
abstract final class AppMotion {
  static const Duration quick = Duration(milliseconds: 140);
  static const Duration standard = Duration(milliseconds: 220);
  static const Duration page = Duration(milliseconds: 280);
  static const Duration background = Duration(milliseconds: 360);

  static const Curve enter = Curves.easeOutCubic;
  static const Curve exit = Curves.easeInCubic;
  static const Curve emphasized = Curves.easeOutQuart;

  /// Short fade + vertical drift used when content changes inside a route.
  static Widget fadeSlideTransition(
    Widget child,
    Animation<double> animation,
  ) {
    final curved = CurvedAnimation(
      parent: animation,
      curve: enter,
      reverseCurve: exit,
    );
    final position = Tween<Offset>(
      begin: const Offset(0, 0.025),
      end: Offset.zero,
    ).animate(curved);
    return FadeTransition(
      opacity: curved,
      child: SlideTransition(position: position, child: child),
    );
  }

  static bool reduceMotion(BuildContext context) =>
      MediaQuery.maybeOf(context)?.disableAnimations ?? false;

  static Duration resolve(BuildContext context, Duration duration) =>
      reduceMotion(context) ? Duration.zero : duration;
}
