# FT-music

**Your music, openly connected.**

FT-music is an open-source music project with a cross-platform client and a separately maintained Go API. The client and backend share a product identity but remain independently testable and deployable.

## Project layout

```text
FT-music/
├── app/       # Flutter client and its Rust/WASM integration
├── brand/     # Canonical product identity and design tokens
└── go/        # Shared Go HTTP API and persistence layer
```

## Backend

The Go service is organized around explicit boundaries:

- `internal/domain` — validated music entities and domain constraints.
- `internal/repository` — context-aware storage contracts.
- `internal/repository/sqlite` — SQLite persistence and schema initialization.
- `internal/repository/memory` — deterministic in-memory implementation for tests.
- `internal/service` — use cases and orchestration.
- `internal/streaming` — provider registration, concurrent search aggregation, a bounded TTL/LRU result cache, in-flight request coalescing, and stream URL validation.
- `internal/server` — versioned JSON HTTP API, request limits, health/readiness checks, and transport middleware.

### Run the API

```bash
cd go
go test ./...
go run ./cmd/ft-music
```

Configuration is environment-based:

| Variable | Default | Purpose |
|---|---|---|
| `FT_MUSIC_ADDR` | `127.0.0.1:8080` | HTTP bind address |
| `FT_MUSIC_DB` | `ft-music.db` | SQLite database path |
| `FT_MUSIC_API_TOKEN` | unset | Bearer token; required for non-loopback binds |

The service exposes `/healthz` for liveness and `/readyz` for storage readiness. It binds to loopback by default. To expose it on a LAN or public interface, set `FT_MUSIC_ADDR` and configure `FT_MUSIC_API_TOKEN`; the server refuses non-loopback binds without a token. The versioned API is documented in [`go/API.md`](go/API.md).

## Streaming and content policy

Provider adapters must use official or otherwise authorized APIs. Each adapter must document its licensing basis, provider terms, attribution requirements, rate limits, and regional restrictions. FT-music does not scrape services, bypass access controls, remove DRM, or distribute music without permission. The project license covers source code only; it does not grant rights to third-party music, artwork, trademarks, or provider services.

## Brand system

FT-music is the product brand. The canonical colors, shape language, logo direction, and typography references are maintained in [`brand/theme.json`](brand/theme.json). Client theme tokens are kept aligned with that specification. Upstream copyright and license notices remain intact where required; upstream identifiers retained for compatibility are internal implementation details, not product branding.

## Licensing

The client includes GPL-licensed upstream-derived source. Preserve its copyright notices, license text, and applicable source-distribution obligations when modifying or distributing the project. See [`LICENSE`](LICENSE), [`NOTICE.md`](NOTICE.md), and the notices retained in the source tree. Review third-party assets and dependency licenses separately before distribution.
