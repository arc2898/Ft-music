# FT-music brand system

FT-music is the product name used in user-facing copy, metadata, API responses, and project documentation.

- **Tagline:** Your music, openly connected
- **Primary accent:** Lime `#A3E635`
- **Supporting lime:** `#B7F36A` and `#D6FFA1`
- **Dark canvas:** `#0B1008`
- **Dark surfaces:** `#151D10` and `#1F2A17`
- **Light canvas:** `#F7FAF2`
- **Light surface:** `#FFFFFF`
- **Typography:** ReThink Sans for interface text; CodePro for technical/metadata accents
- **Mark direction:** A bold lime F monogram with short sound-wave bars on a deep olive-black tile

`theme.json` is the canonical machine-readable token source. Client theme values should map to these tokens rather than introduce independent colors. The editable vector mark is [icon.svg](icon.svg), and the Flutter app uses `app/assets/icons/ft_music_mark.png` for the in-app mark and launcher source.

To regenerate platform launcher assets after changing the raster source, run `dart run flutter_launcher_icons` from the `app/` directory. Keep upstream license/copyright notices and compatibility identifiers where legally or technically required; they are not product branding.
