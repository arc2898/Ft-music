import 'package:flutter/material.dart';

/// FT-music design tokens. Screens should prefer Theme.of(context) so both
/// palettes respond to the persisted appearance setting.
class AppTheme {
  static const primaryTextStyle = TextStyle(fontFamily: 'ReThink-Sans');
  static const secondoryTextStyle = TextStyle(fontFamily: 'ReThink-Sans');
  static const secondoryTextStyleMedium =
      TextStyle(fontFamily: 'ReThink-Sans', fontWeight: FontWeight.w700);
  static const tertiaryTextStyle = TextStyle(fontFamily: 'CodePro');
  static const fontAwesomeRegularFont =
      TextStyle(fontFamily: 'FontAwesome-Regular');
  static const fontAwesomeSolidFont =
      TextStyle(fontFamily: 'FontAwesome-Solids');

  // Dark palette: deep olive-black surfaces with a high-contrast lime accent.
  static const themeColor = Color(0xFF0B1008);
  static const surfaceColor = Color(0xFF151D10);
  static const surfaceElevated = Color(0xFF1F2A17);
  static const surfaceSoft = Color(0xFF2B3A20);
  static const surfaceHighlight = Color(0xFF41552C);
  static const primaryColor1 = Color(0xFFF5F9EE);
  static const primaryColor2 = Color(0xFFB8C7A6);
  static const accentColor1 = Color(0xFFB7F36A);
  static const accentColor1light = Color(0xFFD6FFA1);
  static const accentColor2 = Color(0xFFA3E635);
  static const successColor = Color(0xFF4ADE80);
  static const warningColor = Color(0xFFFBBF24);
  static const errorColor = Color(0xFFFB7185);

  static const cornerSmall = 12.0;
  static const cornerMedium = 18.0;
  static const cornerLarge = 26.0;
  static const motionCurve = Curves.easeOutCubic;
  static const motionFast = Duration(milliseconds: 180);
  static const motionStandard = Duration(milliseconds: 260);

  ThemeData get defaultThemeData => _buildTheme(
        brightness: Brightness.dark,
        background: themeColor,
        surface: surfaceColor,
        elevatedSurface: surfaceElevated,
        softSurface: surfaceSoft,
        outline: surfaceHighlight,
        text: primaryColor1,
        secondaryText: primaryColor2,
        primary: accentColor2,
        secondary: accentColor1,
        onPrimary: const Color(0xFF17220C),
      );

  ThemeData get lightThemeData => _buildTheme(
        brightness: Brightness.light,
        background: const Color(0xFFF7FAF2),
        surface: const Color(0xFFFFFFFF),
        elevatedSurface: const Color(0xFFEEF4E7),
        softSurface: const Color(0xFFE3EBD9),
        outline: const Color(0xFFC7D6B6),
        text: const Color(0xFF1D2916),
        secondaryText: const Color(0xFF58674D),
        primary: const Color(0xFF5C8516),
        secondary: const Color(0xFF76A72A),
        onPrimary: const Color(0xFFFFFFFF),
      );

  ThemeData _buildTheme({
    required Brightness brightness,
    required Color background,
    required Color surface,
    required Color elevatedSurface,
    required Color softSurface,
    required Color outline,
    required Color text,
    required Color secondaryText,
    required Color primary,
    required Color secondary,
    required Color onPrimary,
  }) {
    final scheme = ColorScheme.fromSeed(
      seedColor: primary,
      brightness: brightness,
      surface: surface,
      primary: primary,
      secondary: secondary,
      tertiary: outline,
      onPrimary: onPrimary,
      onSecondary: text,
      onSurface: text,
      error: errorColor,
    );

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      scaffoldBackgroundColor: background,
      dialogBackgroundColor: elevatedSurface,
      canvasColor: background,
      fontFamily: 'ReThink-Sans',
      colorScheme: scheme,
      primaryColor: primary,
      primaryColorDark: primary,
      splashFactory: InkSparkle.splashFactory,
      visualDensity: VisualDensity.adaptivePlatformDensity,
      iconTheme: IconThemeData(color: text, size: 22),
      appBarTheme: AppBarTheme(
        backgroundColor: Colors.transparent,
        foregroundColor: text,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        centerTitle: false,
        titleTextStyle: TextStyle(
          color: text,
          fontFamily: 'ReThink-Sans',
          fontSize: 22,
          fontWeight: FontWeight.w800,
          letterSpacing: -0.4,
        ),
      ),
      cardTheme: CardThemeData(
        color: surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(cornerMedium)),
        ),
      ),
      listTileTheme: ListTileThemeData(
        dense: true,
        minLeadingWidth: 0,
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 5),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(cornerMedium),
        ),
        iconColor: secondaryText,
        textColor: text,
        subtitleTextStyle: secondoryTextStyle.copyWith(
          color: secondaryText,
          fontSize: 12,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: primary,
          foregroundColor: onPrimary,
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(cornerMedium),
          ),
          textStyle: secondoryTextStyleMedium.copyWith(fontSize: 14),
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: primary,
          foregroundColor: onPrimary,
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(cornerMedium),
          ),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: scheme.primary,
          side: BorderSide(color: outline),
          padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 13),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(cornerMedium),
          ),
        ),
      ),
      chipTheme: ChipThemeData(
        backgroundColor: elevatedSurface,
        selectedColor: primary,
        disabledColor: softSurface,
        labelStyle: secondoryTextStyleMedium.copyWith(
          color: text,
          fontSize: 12,
        ),
        secondaryLabelStyle: secondoryTextStyleMedium.copyWith(
          color: onPrimary,
        ),
        side: BorderSide.none,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(999),
        ),
      ),
      bottomSheetTheme: BottomSheetThemeData(
        backgroundColor: elevatedSurface,
        modalBackgroundColor: elevatedSurface,
        surfaceTintColor: Colors.transparent,
        showDragHandle: true,
        dragHandleColor: outline,
        shape: const RoundedRectangleBorder(
          borderRadius:
              BorderRadius.vertical(top: Radius.circular(cornerLarge)),
        ),
      ),
      dialogTheme: DialogThemeData(
        backgroundColor: elevatedSurface,
        surfaceTintColor: Colors.transparent,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(cornerLarge)),
        ),
      ),
      pageTransitionsTheme: const PageTransitionsTheme(
        builders: {
          TargetPlatform.android: FadeForwardsPageTransitionsBuilder(),
          TargetPlatform.iOS: FadeForwardsPageTransitionsBuilder(),
          TargetPlatform.linux: FadeForwardsPageTransitionsBuilder(),
          TargetPlatform.macOS: FadeForwardsPageTransitionsBuilder(),
          TargetPlatform.windows: FadeForwardsPageTransitionsBuilder(),
        },
      ),
      scrollbarTheme: ScrollbarThemeData(
        thumbColor: WidgetStatePropertyAll(primary.withValues(alpha: .72)),
        trackColor: WidgetStatePropertyAll(softSurface.withValues(alpha: .55)),
        interactive: true,
        radius: const Radius.circular(12),
        thickness: const WidgetStatePropertyAll(4),
      ),
      progressIndicatorTheme: ProgressIndicatorThemeData(
        color: primary,
        linearTrackColor: softSurface,
        circularTrackColor: softSurface,
      ),
      textSelectionTheme: TextSelectionThemeData(
        cursorColor: primary,
        selectionColor: primary.withValues(alpha: .28),
        selectionHandleColor: primary,
      ),
      switchTheme: SwitchThemeData(
        thumbColor: WidgetStatePropertyAll(onPrimary),
        trackColor: WidgetStateProperty.resolveWith((states) =>
            states.contains(WidgetState.selected) ? primary : softSurface),
        trackOutlineColor: const WidgetStatePropertyAll(Colors.transparent),
      ),
      searchBarTheme: SearchBarThemeData(
        backgroundColor: WidgetStatePropertyAll(elevatedSurface),
        surfaceTintColor: const WidgetStatePropertyAll(Colors.transparent),
        elevation: const WidgetStatePropertyAll(0),
        padding: const WidgetStatePropertyAll(
          EdgeInsets.symmetric(horizontal: 18),
        ),
        shape: WidgetStatePropertyAll(
          RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(cornerLarge),
            side: BorderSide(color: outline),
          ),
        ),
        hintStyle: WidgetStatePropertyAll(TextStyle(color: secondaryText)),
      ),
      popupMenuTheme: PopupMenuThemeData(
        color: elevatedSurface,
        surfaceTintColor: Colors.transparent,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(cornerMedium)),
        ),
        textStyle: TextStyle(color: text, fontFamily: 'ReThink-Sans'),
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: elevatedSurface,
        indicatorColor: primary,
        surfaceTintColor: Colors.transparent,
        height: 72,
        labelTextStyle: const WidgetStatePropertyAll(
          TextStyle(fontFamily: 'ReThink-Sans', fontWeight: FontWeight.w700),
        ),
      ),
    );
  }
}

// Backward-compatible alias used throughout the existing client.
// ignore: camel_case_types
typedef Default_Theme = AppTheme;
