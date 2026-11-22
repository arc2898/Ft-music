#!/usr/bin/env python3
"""Regenerate every checked-in FT-music launcher icon from one canonical PNG.

Run from any working directory with Pillow installed:
    python tooling/create_ft_music_logo.py
or from app/: 
    dart run flutter_launcher_icons

The Flutter tool is the preferred cross-platform generator. This script also
keeps the legacy hand-maintained icon sets in sync.
"""
from __future__ import annotations

import json
import re
from pathlib import Path
from PIL import Image

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "assets" / "icons" / "ft_music_mark.png"
RES = ROOT / "android" / "app" / "src" / "main" / "res"

if not SOURCE.is_file():
    raise SystemExit(f"Missing canonical FT-music icon: {SOURCE}")

icon = Image.open(SOURCE).convert("RGBA")


def save_square(path: Path, size: int) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    icon.resize((size, size), Image.Resampling.LANCZOS).save(path, optimize=True)


# Android legacy launcher bitmaps.
for density, size in {
    "mdpi": 48,
    "hdpi": 72,
    "xhdpi": 96,
    "xxhdpi": 144,
    "xxxhdpi": 192,
}.items():
    for name in ("ic_launcher.png", "launcher_icon.png"):
        save_square(RES / f"mipmap-{density}" / name, size)

# iOS asset catalog: use the filenames already declared by Xcode.
ios_set = ROOT / "ios" / "Runner" / "Assets.xcassets" / "AppIcon.appiconset"
contents = ios_set / "Contents.json"
if contents.is_file():
    for item in json.loads(contents.read_text(encoding="utf-8")).get("images", []):
        filename = item.get("filename")
        if not filename:
            continue
        match = re.match(r"Icon-App-(\d+(?:\.5)?)[xX]\d+(?:@([123])x)?\.png$", filename)
        if match:
            logical = float(match.group(1))
            scale = int(match.group(2) or 1)
            save_square(ios_set / filename, round(logical * scale))

# macOS icon set.
mac_set = ROOT / "macos" / "Runner" / "Assets.xcassets" / "AppIcon.appiconset"
for size in (16, 32, 64, 128, 256, 512, 1024):
    save_square(mac_set / f"app_icon_{size}.png", size)

# Windows multi-resolution icon.
windows_icon = ROOT / "windows" / "runner" / "resources" / "app_icon.ico"
windows_icon.parent.mkdir(parents=True, exist_ok=True)
icon.save(
    windows_icon,
    format="ICO",
    sizes=[(16, 16), (24, 24), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)],
)

# Store listing icon.
save_square(ROOT / "fastlane" / "metadata" / "android" / "en-US" / "images" / "icon.png", 512)

print("Regenerated FT-music launcher icons for Android, iOS, macOS, Windows and store metadata.")
