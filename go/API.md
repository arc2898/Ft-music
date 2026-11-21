# FT-music API v1

Base path: `/api/v1`  
Content type: `application/json; charset=utf-8`

The API is intentionally provider-neutral. The client talks to the FT-music contract; provider-specific behavior stays behind the `streaming.Provider` interface.

## Authentication

The service binds to `127.0.0.1:8080` by default and does not require a token for local-only use. When `FT_MUSIC_API_TOKEN` is configured, all API endpoints except liveness/readiness (`/healthz`, `/readyz`, and `/api/v1/health`) require `Authorization: Bearer <token>`. A non-loopback bind requires a non-empty token; the process refuses to start otherwise. Keep the token secret and terminate public traffic over HTTPS at a trusted reverse proxy.

## Health

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/healthz` | Process liveness |
| `GET` | `/readyz` | Storage readiness |
| `GET` | `/api/v1/health` | Backwards-compatible API health response |

`/healthz` is suitable for liveness probes. `/readyz` returns `503` when the storage layer cannot be reached.

## Library

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/v1/library` | List or search saved tracks |
| `GET` | `/api/v1/library/{id}` | Retrieve one saved track |
| `POST` | `/api/v1/library` | Add or update a saved track by ID |
| `DELETE` | `/api/v1/library/{id}` | Remove a saved track |

`GET /library` accepts:

- `q` — optional case-insensitive match against title, artist, or album (maximum 256 characters).
- `limit` — optional integer from 1 to 100; defaults to 50.
- `offset` — optional non-negative integer; defaults to 0.

The response remains a JSON array for simple client compatibility. Pagination metadata is supplied in `X-Total-Count`, `X-Limit`, and `X-Offset` response headers. Results have deterministic creation-time and ID ordering.

`POST /library` accepts one JSON object and rejects unknown properties, trailing JSON values, invalid content types, bodies larger than 1 MiB, missing IDs/titles/artists, negative durations, and non-HTTP(S) artwork URLs. Existing IDs are upserted.

Example:

```json
{
  "id": "provider:track-123",
  "title": "Signal",
  "artist": "FT-music",
  "album": "Night Drive",
  "duration_seconds": 215,
  "artwork_url": "https://cdn.example.test/cover.jpg"
}
```

## Providers

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/v1/providers` | List registered provider names |
| `GET` | `/api/v1/search?q={query}` | Search authorized remote providers |
| `GET` | `/api/v1/providers/{provider}/tracks/{id}/stream` | Resolve a provider stream URL |

Remote search fans out to registered providers, merges duplicate title/artist pairs, and returns partial successful results when another provider fails. Successful normalized queries are cached in a bounded in-memory LRU for two minutes; concurrent identical requests share one provider search. With no providers configured it returns `503`. Provider stream URLs must be absolute HTTP(S) URLs without embedded user information.

No provider is enabled by default. Integrations must be explicitly registered and must comply with the provider's API and content terms.

## Error format

All application errors use a consistent envelope:

```json
{
  "error": {
    "code": "track_not_found",
    "message": "track was not found"
  }
}
```

The server also returns an `X-Request-ID` header for request correlation. Responses include `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, and `Cache-Control: no-store`.

## Operational behavior

The HTTP server applies read-header, read, write, and idle timeouts. It shuts down gracefully on `SIGINT` and `SIGTERM`. SQLite uses a busy timeout, WAL mode for file-backed databases, and a single connection to keep write behavior predictable.
