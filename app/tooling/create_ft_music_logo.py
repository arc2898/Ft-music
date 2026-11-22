#!/usr/bin/env python3
"""Regenerate legacy Android launcher bitmaps from the canonical FT-music icon.

The canonical artwork is assets/icons/ft_music_mark.png. Do not draw a second
logo here: Flutter Launcher Icons reads the same source from pubspec.yaml for
Android adaptive icons, iOS, macOS and Windows.
"""
from pathlib import Path
from PIL import Image

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "assets" / "icons" / "ft_music_mark.png"
RES = ROOT / "android" / "app" / "src" / "main" / "res"

if not SOURCE.is_file():
    raise SystemExit(f"Missing canonical FT-music icon: {SOURCE}")

icon = Image.open(SOURCE).convert("RGBA")
for density, size in {
    "mdpi": 48,
    "hdpi": 72,
    "xhdpi": 96,
    "xxhdpi": 144,
    "xxxhdpi": 192,
}.items():
    target = RES / f"mipmap-{density}"
    target.mkdir(parents=True, exist_ok=True)
    resized = icon.resize((size, size), Image.Resampling.LANCZOS)
    resized.save(target / "ic_launcher.png", optimize=True)
    resized.save(target / "launcher_icon.png", optimize=True)

print("Regenerated legacy Android launcher icons from ft_music_mark.png")
print("For all platform icon sets, run: dart run flutter_launcher_icons")
