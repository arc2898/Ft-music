# FT-music splash and launcher branding

The startup and plugin-bootstrap surfaces use the canonical FT-music lime mark
from `assets/icons/ft_music_mark.png`. The scalable companion is
`assets/icons/ft_music_mark.svg`.

## Plugin bootstrap overlay

- Uses the bundled Unageo typeface explicitly and disables text decoration so
  Android accessibility/font settings cannot make the labels look like links.
- Shows the current plugin-bootstrap percentage in both the ring and a slim
  linear progress bar.
- Uses the same lime accent as the app theme and the new brand mark.
- Keeps retry, offline, and failure handling intact.

## Launcher assets

The canonical source is configured in `pubspec.yaml` for Android, iOS, macOS,
and Windows. Android's adaptive icon uses the matching dark background and
lime vector foreground. The legacy Android density launcher bitmaps are
regenerated from the canonical source; stale density-specific foreground and
background PNG layers and the unused `ft_music_lighter.png` duplicate were
removed.

After changing the source mark, regenerate the platform sets from `app/`:

```sh
dart run flutter_launcher_icons
python tooling/create_ft_music_logo.py
```

## Validation

From `app/`:

```sh
flutter analyze
flutter test
```

Also install the Android build on a device and verify the launcher icon on
pre-Android-8 and adaptive-icon launchers, then check the startup splash and
plugin bootstrap at multiple display sizes and with large system font scaling.
