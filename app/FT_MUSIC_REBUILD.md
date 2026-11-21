# FT-music implementation notes

FT-music is maintained as a branded product with two independently testable parts: the Flutter client in `app/` and the shared Go API in `go/`. The client source is derived from a GPL-licensed upstream Flutter/Rust project; its required copyright and license notices must remain intact. Product-facing copy, metadata, documentation, and design tokens use FT-music branding.

## Backend baseline

The Go API is the shared service boundary for future clients. It is deliberately separate from the Flutter build and can be developed, tested, and deployed on its own.

Current backend capabilities:

- validated, normalized track metadata with consistent domain constraints;
- context-aware memory and SQLite repositories;
- SQLite initialization, WAL mode for file databases, busy timeout, and readiness checks;
- a provider registry with isolated provider failures, parallel search aggregation, de-duplication, and stream URL validation;
- versioned JSON routes for health, readiness, library CRUD, remote search, provider discovery, and stream resolution;
- bounded request bodies, strict JSON decoding, pagination, request IDs, security headers, API-key protection, and graceful shutdown.

No streaming provider is enabled by default. A provider adapter must be explicitly registered and must use an authorized API with documented licensing, attribution, rate limits, and regional restrictions.

## Run and test the backend

```bash
cd go
go test ./...
go run ./cmd/ft-music
```

Configuration:

| Variable | Default | Purpose |
|---|---|---|
| `FT_MUSIC_ADDR` | `127.0.0.1:8080` | HTTP bind address |
| `FT_MUSIC_DB` | `ft-music.db` | SQLite database path |
| `FT_MUSIC_API_TOKEN` | unset | Bearer token; required for non-loopback binds |

See [`../go/API.md`](../go/API.md) for endpoint contracts and operational behavior.

## Branding and source lineage

[`../brand/theme.json`](../brand/theme.json) is the canonical FT-music brand token file. The Flutter theme follows its violet/cyan palette, typography, shape, and motion tokens. Keep FT-music as the public product identity. Preserve upstream copyright notices, license files, and internal compatibility identifiers required by generated Flutter/Rust bridge code; do not present upstream names as FT-music product branding.
