# FT-music brand system

FT-music is the product name used in user-facing copy, metadata, API responses, and project documentation.

- **Tagline:** Your music, openly connected
- **Primary:** Violet `#7C3AED`
- **Accent:** Cyan `#22D3EE`
- **Canvas:** Midnight `#0B1020`
- **Surfaces:** `#151B2E` and `#202943`
- **Typography:** Unageo for interface text; CodePro for technical/metadata accents
- **Mark direction:** Two overlapping sound-wave arcs forming an F

`theme.json` is the canonical machine-readable token source. Client theme values should map to these tokens rather than introduce independent colors. Keep upstream license/copyright notices and compatibility identifiers where legally or technically required; they are not product branding.


## Lime refresh

- Primary accent: `#A3E635`; supporting lime: `#B7F36A`.
- Dark surfaces use olive-black neutrals; the companion light palette is defined in [theme.json](theme.json).
- Interface typography: ReThink Sans; technical typography: Code Pro.
- The source app mark is [icon.svg](icon.svg). The Flutter launcher generator reads `app/assets/icons/ft_music_mark.png`; regenerate platform launcher assets with `dart run flutter_launcher_icons` after changing that raster source.
