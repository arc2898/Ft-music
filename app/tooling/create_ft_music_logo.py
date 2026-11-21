from pathlib import Path
from PIL import Image, ImageDraw

SIZE = 1024
BACKGROUND = (11, 16, 8, 255)
TILE = (15, 23, 10, 255)
LIME = (190, 242, 66, 255)
LIME_DARK = (151, 210, 44, 255)
RING = (72, 105, 27, 255)

image = Image.new("RGBA", (SIZE, SIZE), BACKGROUND)
draw = ImageDraw.Draw(image)

# Rounded tile, inner keyline, and quiet circular sound-wave ring.
draw.rounded_rectangle((18, 18, SIZE - 18, SIZE - 18), radius=224, fill=TILE,
                       outline=(54, 74, 28, 255), width=10)
draw.ellipse((180, 180, 844, 844), outline=RING, width=16)

# A compact F monogram with a vertical stem and two broad bars.
draw.rounded_rectangle((296, 264, 408, 760), radius=56, fill=LIME)
draw.rounded_rectangle((356, 264, 736, 376), radius=56, fill=LIME)
draw.rounded_rectangle((356, 456, 656, 568), radius=56, fill=LIME)

# Three sound-wave beats complete the mark without adding fine detail.
draw.rounded_rectangle((692, 584, 768, 660), radius=34, fill=LIME_DARK)
draw.rounded_rectangle((664, 680, 768, 756), radius=34, fill=LIME)
draw.rounded_rectangle((600, 776, 768, 852), radius=34, fill=LIME_DARK)
draw.ellipse((724, 192, 800, 268), fill=LIME)

assets = Path("assets/icons")
assets.mkdir(parents=True, exist_ok=True)
image.save(assets / "ft_music_mark.png", optimize=True)
image.resize((512, 512), Image.Resampling.LANCZOS).save(
    assets / "ft_music_lighter.png", optimize=True
)

# Regenerate the legacy Android launcher bitmap sizes from the same source.
res = Path("android/app/src/main/res")
for density, size in {
    "mdpi": 48,
    "hdpi": 72,
    "xhdpi": 96,
    "xxhdpi": 144,
    "xxxhdpi": 192,
}.items():
    resized = image.resize((size, size), Image.Resampling.LANCZOS)
    target = res / f"mipmap-{density}"
    target.mkdir(parents=True, exist_ok=True)
    resized.save(target / "ic_launcher.png", optimize=True)
    resized.save(target / "launcher_icon.png", optimize=True)
