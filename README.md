# ClipBridgeServer

A lightweight self-hosted clipboard stack service for syncing text, pairing
devices, managing history, and serving an embedded Web UI from one portable
server binary.

## Current Stage

This repository is currently in phase 7, with phase 6.5 included:

- phase 6.5: Web UI Quick Clipboard Panel
- phase 7: cleaner plus retention policy

The service can now do three important things at the same time:

- act as the central clipboard API and history store
- let the browser act as a lightweight manual clipboard client
- run as a long-lived self-hosted service with automatic retention controls

Deployment is still intentionally simple:

- one server binary
- one `config.yaml`
- one `data/` directory

## Progress So Far

Delivered through phase 7:

- `go run ./cmd/server -config config.yaml` starts the HTTP server
- `internal/config` loads and validates YAML config
- `internal/store` opens SQLite, creates `data/clipbridge.db`, and runs migrations
- `GET /api/health` works without a token
- bearer token auth protects every other API route
- one-time pairing codes can mint long-lived `device_token` values
- device tokens can call clipboard APIs
- `POST /api/clipboard/text` stores text clipboard items
- `GET /api/clipboard/latest` returns the newest text item
- `GET /api/clipboard/history` returns history ordered newest first
- `GET /api/clipboard/items/{id}` returns one record
- `DELETE /api/clipboard/items/{id}` soft-deletes one record
- `POST /api/clipboard/items/{id}/favorite` favorites one record
- `DELETE /api/clipboard/items/{id}/favorite` removes one favorite
- `GET /api/favorites` returns favorite clipboard items
- built-in categories `text`, `image`, `link`, and `file` are seeded automatically
- `GET /api/categories` lists built-in and custom categories
- `POST /api/categories` creates one custom category
- `PATCH /api/clipboard/items/{id}/category` updates one item's category
- `GET /api/clipboard/history?category=text` filters history by category
- `POST /api/admin/cleanup/run` triggers one manual cleanup pass
- `GET /api/admin/cleanup/status` returns the latest cleanup summary
- `GET /api/admin/storage/status` returns current storage usage
- `GET /api/settings/cleanup` returns the persisted cleanup policy
- `PATCH /api/settings/cleanup` updates the cleanup policy without restarting the server
- cleanup policy is stored in the `settings` table and reloaded after restart
- favorites are skipped by automatic cleanup
- the background cleaner can enforce TTL, max history count, and max total storage size
- successful API responses use a top-level `data`
- failed API responses use a top-level `error`
- request size limits, text size limits, CORS, and request logs are enabled
- `GET /` serves an embedded Web UI from the Go binary
- the Web UI can upload browser text to the shared clipboard stack through the existing clipboard API
- the Web UI can read server latest, copy latest into the browser clipboard, and refresh latest manually
- the Web UI can show history, favorites, devices, pairing codes, cleanup status, and cleanup policy

Not implemented yet:

- image clipboard upload and download
- file clipboard upload and download
- automatic desktop or mobile background clipboard watching
- WebDAV sync
- browser session login flow separate from bearer tokens

## Quick Start

1. Copy the example config:

```bash
cp configs/config.example.yaml config.yaml
```

2. Start the server:

```bash
go run ./cmd/server -config config.yaml
```

3. Verify the health endpoint:

```bash
curl http://127.0.0.1:8787/api/health
```

Expected response:

```json
{"data":{"ok":true,"version":"0.1.0"}}
```

4. Open the embedded Web UI:

```text
http://127.0.0.1:8787/
```

After startup, the server will create:

- `data/clipbridge.db`

## How To Use

### 1. Use The Admin Token

By default the example config uses:

```text
dev-token-123
```

Every protected API request must include:

```http
Authorization: Bearer <token>
```

### 2. Upload Text From A Client Or Browser

Classic text upload:

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/text \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"text":"hello from ClipBridge"}'
```

Web UI style upload with source metadata, using the same API:

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/text \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"content":"from browser","source_device_id":"web-ui","source_device_name":"Web UI"}'
```

### 3. Read Latest And History

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/latest

curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/history
```

Filter by category:

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  'http://127.0.0.1:8787/api/clipboard/history?category=text'
```

### 4. Favorite Important Records

```bash
curl -X POST \
  -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/items/1/favorite

curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/favorites

curl -X DELETE \
  -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/items/1/favorite
```

### 5. Organize With Categories

List categories:

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/categories
```

Create a custom category:

```bash
curl -X POST http://127.0.0.1:8787/api/categories \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"name":"work"}'
```

Move one record into a category:

```bash
curl -X PATCH http://127.0.0.1:8787/api/clipboard/items/1/category \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"category":"work"}'
```

### 6. Pair A Device

Generate a one-time pairing code:

```bash
curl -X POST http://127.0.0.1:8787/api/auth/pairing-codes \
  -H 'Authorization: Bearer dev-token-123'
```

Exchange it on the client:

```bash
curl -X POST http://127.0.0.1:8787/api/auth/pair \
  -H 'Content-Type: application/json' \
  -d '{"pairing_code":"ABCDEFGH","device_name":"My Laptop"}'
```

After that, the returned `device_token` can call clipboard APIs without asking
the user to copy the long admin token manually.

### 7. Inspect And Run Cleanup

Check current storage usage:

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/admin/storage/status
```

Check the current cleanup policy:

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/settings/cleanup
```

Update the cleanup policy:

```bash
curl -X PATCH http://127.0.0.1:8787/api/settings/cleanup \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"ttl_hours":168,"max_items":1000,"max_total_size_mb":2048,"interval_minutes":30,"enabled":true}'
```

Run cleanup manually:

```bash
curl -X POST http://127.0.0.1:8787/api/admin/cleanup/run \
  -H 'Authorization: Bearer dev-token-123'
```

## Web UI Usage

Open `http://127.0.0.1:8787/` and:

1. paste either the admin token or a device token into the token box
2. press `Save`
3. use `Quick Clipboard` to type text and upload it into the shared clipboard stack
4. use `Copy Server Latest` to write the current server latest text into the browser clipboard
5. use `Favorites` and `Clipboard History` to organize records
6. use `Retention Status` and `Cleanup Policy` with the admin token when operating the server long-term
7. use `Pair Devices` and `Paired Devices` with the admin token when onboarding clients

Notes:

- admin token can access pairing, device management, cleanup status, and cleanup settings APIs
- device token can access clipboard, favorites, and category APIs
- the browser clipboard write path is manual only and triggered by a user click
- the token is stored only in the current browser's local storage

## Cleanup Behavior

The current retention rules are:

- non-favorite items older than the configured TTL are deleted automatically
- if active history count exceeds `max_items`, the oldest non-favorite items are deleted first
- if total storage exceeds `max_total_size_mb`, the oldest non-favorite file items are deleted first, then other non-favorites if needed
- favorite items are never deleted by automatic cleanup
- manual deletion still works even when an item is favorited

## Directory Layout

- `cmd/server/main.go`: application entrypoint; reads config, initializes SQLite, starts cleanup service, creates router, starts HTTP server
- `internal/config/config.go`: config structures, YAML loading, validation, and server address helper
- `internal/config/defaults.go`: default config values for the current stage
- `internal/auth/middleware.go`: bearer token parsing, pairing code generation, and secret hashing helpers
- `internal/cleanup/service.go`: background cleaner, manual cleanup runner, and cleanup status assembly
- `internal/store/sqlite.go`: SQLite connection lifecycle plus clipboard, favorites, categories, settings, and cleanup storage methods
- `internal/store/settings.go`: cleanup policy validation helpers
- `internal/store/device_pairing.go`: pairing code, device token, device list, and device revoke storage logic
- `internal/store/migrations.go`: migration runner and embedded SQL migration loading
- `internal/store/sqlite_test.go`: regression test for database and table creation
- `internal/store/text_items_test.go`: regression test for clipboard CRUD, metadata, favorites, and categories
- `internal/store/device_pairing_test.go`: storage tests for pairing and device revocation
- `migrations/001_init.sql`: initial schema for base tables
- `migrations/002_device_pairing.sql`: device and pairing code schema
- `migrations/003_favorites_categories.sql`: built-in category seeding and legacy category backfill
- `migrations/004_cleanup_metadata.sql`: size and expiration columns for retention logic
- `internal/api/router.go`: route registration, auth middleware, CORS, request size protection, and logging
- `internal/api/health_handler.go`: health check handler
- `internal/api/auth_handler.go`: pairing code and device management handlers
- `internal/api/clipboard_handler.go`: clipboard, favorites, and category HTTP handlers
- `internal/api/cleanup_handler.go`: manual cleanup, storage status, and cleanup settings handlers
- `internal/api/response.go`: shared success and error response helpers
- `internal/api/clipboard_handler_test.go`: HTTP-level tests for auth, pairing, quick clipboard, favorites, categories, and cleanup flows
- `web/embedded.go`: `go:embed` entrypoint for the built frontend bundle
- `web/dist/index.html`: embedded Web UI entry document
- `web/dist/app.css`: embedded Web UI styles
- `web/dist/app.js`: embedded Web UI behavior
- `configs/config.example.yaml`: example config for local startup
- `docs/api.md`: HTTP API reference for the current phase
- `docs/config.md`: active config keys used by the implementation
- `docs/roadmap.md`: staged build plan for later phases
- `CHANGELOG.md`: phase-by-phase change record

## Notes

The Go module path is normalized to:

`github.com/cinmou/ClipBridgeServer`

All internal imports should use that exact casing so builds stay consistent
across macOS, Windows, Linux, and OpenWrt.
