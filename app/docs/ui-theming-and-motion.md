# UI theming and motion

FT-music screens use the active Material color scheme for foreground, surface,
outline, and accent colors. In particular, Local Music, Library, Search,
shared song rows, and library item rows must not use the legacy dark-only
`Default_Theme` palette for screen backgrounds or primary text.

## Motion contract

Use `AppMotion` tokens for transitions:

- `quick` for immediate feedback such as selection and clear actions.
- `standard` for list items, text, and state transitions.
- `page` for top-level navigation.
- `background` for low-frequency background changes.

Resolve each duration with `AppMotion.resolve(context, duration)`. This
returns zero when the platform requests animations to be disabled. Prefer
`easeOutCubic` for entering motion and `easeInCubic` for exiting
motion. Avoid elastic/bounce effects, large scaling, long stagger chains, and
large animated repaint regions for ordinary navigation.

List stagger is capped to a small number of early rows; long lists should not
delay later items. Keep layout and data loading independent of animation state.
