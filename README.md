# ClipBridgeServer

A lightweight self-hosted clipboard stack service for syncing text, pairing
devices, managing history, and serving an embedded Web UI from one portable
server binary.

Language:

- English: `README.md`
- 简体中文: `README.zh-CN.md`

## Beta 1

This README now serves as the first formal Beta document for
ClipBridgeServer.

Beta 1 means:

- the server is already usable as a self-hosted clipboard hub
- the deployment model is stable enough for normal users to try
- the embedded Web UI is good enough for daily management
- some important future work is still intentionally marked as preview or not yet shipped

## Current Stage

This repository is currently in phase 11, with WebDAV sync preview now added on
top of deployment support:

- phase 8: image, file, and link support
- phase 9: settings API plus Web settings
- phase 10: deployment and release packaging
- phase 11: WebDAV sync preview

The service can now do three important things at the same time:

- act as the central clipboard API and history store
- let the browser act as a lightweight manual clipboard client
- run as a long-lived self-hosted service with automatic retention controls

At this point, the project is already usable as a self-hosted clipboard server
for:

- text, link, image, and small file history
- device pairing with long-lived device tokens
- favorites, categories, and retention policy
- embedded Web UI management
- manual WebDAV sync preview

Deployment is still intentionally simple:

- one server binary
- one `config.yaml`
- one `data/` directory

This phase adds:

- local multi-platform release builds through `scripts/build-release.sh`
- GitHub Actions release automation
- Docker and Docker Compose examples
- `systemd`, `launchd`, Windows NSSM, and OpenWrt service examples
- Homebrew Tap and Scoop manifest templates
- manual WebDAV sync preview with persisted settings and sync status

## Progress So Far

Delivered through phase 9:

- `go run ./cmd/server -config config.yaml` starts the HTTP server
- `internal/config` loads and validates YAML config
- `internal/store` opens SQLite, creates `data/clipbridge.db`, and runs migrations
- `GET /api/health` works without a token
- bearer token auth protects every other API route
- one-time pairing codes can mint long-lived `device_token` values
- device tokens can call clipboard APIs
- `POST /api/clipboard/text` stores text clipboard items
- `POST /api/clipboard/link` stores one link clipboard item
- `POST /api/clipboard/file` stores one image or file clipboard item through `multipart/form-data`
- `GET /api/clipboard/latest` returns the newest clipboard item
- `GET /api/clipboard/history` returns history ordered newest first
- `GET /api/clipboard/items/{id}` returns one record
- `GET /api/clipboard/items/{id}/file` streams one stored image or file back to the client
- `DELETE /api/clipboard/items/{id}` soft-deletes one record
- `POST /api/clipboard/items/{id}/favorite` favorites one record
- `DELETE /api/clipboard/items/{id}/favorite` removes one favorite
- `GET /api/favorites` returns favorite clipboard items
- built-in categories `text`, `image`, `link`, and `file` are seeded automatically
- history can now distinguish and return `text`, `image`, `link`, and `file` items
- SQLite stores file metadata only; uploaded bytes stay on disk under the data directory
- uploaded images can be previewed and downloaded through the existing item API
- uploaded files are streamed to disk while computing SHA-256 and MIME type
- `GET /api/categories` lists built-in and custom categories
- `POST /api/categories` creates one custom category
- `PATCH /api/clipboard/items/{id}/category` updates one item's category
- `GET /api/clipboard/history?category=text` filters history by category
- `POST /api/admin/cleanup/run` triggers one manual cleanup pass
- `GET /api/admin/cleanup/status` returns the latest cleanup summary
- `GET /api/admin/storage/status` returns current storage usage
- `GET /api/settings/cleanup` returns the persisted cleanup policy
- `PATCH /api/settings/cleanup` updates the cleanup policy without restarting the server
- `GET /api/settings` returns the combined runtime settings and startup-only settings
- `PATCH /api/settings` can change the admin token at runtime
- `GET /api/settings/limits` returns the persisted runtime upload limits
- `PATCH /api/settings/limits` updates text, image, file, link, and request size limits without restart
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

Delivered in phase 10:

- `go build ./cmd/server` produces the standalone server binary
- `scripts/build-release.sh` cross-builds the planned release targets into `dist/`
- `.github/workflows/release.yml` runs tests and publishes tagged release assets
- `Dockerfile` packages the server as a single-container service
- `docker-compose.example.yml` provides a one-command local container example
- `deploy/systemd/clipbridge-server.service` provides a Linux boot-time service example
- `deploy/launchd/com.cinmou.clipbridge-server.plist` provides a macOS auto-start example
- `deploy/windows/nssm-install.ps1` provides a Windows long-running service example
- `deploy/openwrt/clipbridge-server.init` provides an OpenWrt init script example
- `deploy/homebrew/clipbridge-server.rb` and `deploy/scoop/clipbridge-server.json` provide package manager templates for later publishing

Delivered in phase 11:

- `GET /api/settings/webdav` and `PATCH /api/settings/webdav` persist WebDAV backend settings
- `POST /api/admin/webdav/test` validates the stored WebDAV connection settings
- `POST /api/admin/webdav/sync` runs one manual WebDAV sync pass
- `GET /api/admin/webdav/status` returns the latest test and sync status
- clipboard items now get a deterministic sync key for import/export deduplication
- WebDAV sync pushes `manifest.json`, `items/*.json`, and `files/*.bin`
- remote text, image, link, and file items can be imported back into local history
- the embedded Web UI can configure WebDAV, test the connection, and trigger sync

Not implemented yet:

- automatic desktop or mobile background clipboard watching
- automatic background WebDAV sync
- browser session login flow separate from bearer tokens
- deep conflict resolution across multiple writers

## What Comes Next

The main work planned after this Beta is:

- desktop clients that can watch the local clipboard in the background
- tray or menu bar resident mode for Windows, Linux, and macOS
- automatic background WebDAV sync instead of manual sync only
- stronger sync conflict handling and clearer sync logs
- optional end-to-end encryption so the server cannot read clipboard contents
- richer browser and client login flows beyond raw bearer tokens
- better import and export tooling for long-term self-hosted backups

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

## Deployment

The shortest path for normal users is now:

1. download the matching binary from GitHub Releases
2. place it beside `config.yaml`
3. create an empty `data/` directory
4. run `./clipbridge-server -config ./config.yaml`

For packaged deployment examples, see:

- `docs/deployment.md`
- `docs/deployment.zh-CN.md`
- `README.zh-CN.md`

## Build A Release

If you just want one binary for your current machine:

```bash
go build -o clipbridge-server ./cmd/server
```

If you want the full Beta release set:

```bash
bash scripts/build-release.sh
```

That writes these files into `dist/`:

- `clipbridge-server-linux-amd64`
- `clipbridge-server-linux-arm64`
- `clipbridge-server-darwin-amd64`
- `clipbridge-server-darwin-arm64`
- `clipbridge-server-windows-amd64.exe`

Recommended local release flow:

1. run `env GOCACHE=$(pwd)/.gocache go test ./...`
2. run `bash scripts/build-release.sh`
3. open `dist/`
4. pick the binary that matches the target platform

If you want GitHub Releases assets instead of local files:

1. commit your changes
2. create a version tag such as `v0.11.0-beta1`
3. push the tag
4. let `.github/workflows/release.yml` build and upload the artifacts

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

### 3.5. Upload Links And Files

Upload a link:

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/link \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","source_device_id":"web-ui","source_device_name":"Web UI"}'
```

Upload an image or file:

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/file \
  -H 'Authorization: Bearer dev-token-123' \
  -F 'file=@./demo.png' \
  -F 'source_device_id=web-ui' \
  -F 'source_device_name=Web UI'
```

Download one stored file:

```bash
curl -L -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/items/1/file
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

### 8. Change Runtime Limits And Token

Read combined settings:

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/settings
```

Change the admin token:

```bash
curl -X PATCH http://127.0.0.1:8787/api/settings \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"admin_token":"new-admin-token"}'
```

Update runtime limits:

```bash
curl -X PATCH http://127.0.0.1:8787/api/settings/limits \
  -H 'Authorization: Bearer new-admin-token' \
  -H 'Content-Type: application/json' \
  -d '{"min_text_bytes":1,"max_text_bytes":262144,"min_image_bytes":1,"max_image_bytes":10485760,"min_file_bytes":1,"max_file_bytes":52428800,"min_link_bytes":1,"max_link_bytes":8192,"max_request_bytes":62914560}'
```

### 9. Preview WebDAV Sync

Save WebDAV settings:

```bash
curl -X PATCH http://127.0.0.1:8787/api/settings/webdav \
  -H 'Authorization: Bearer new-admin-token' \
  -H 'Content-Type: application/json' \
  -d '{"enabled":true,"url":"https://dav.example.com/remote.php/dav/files/user","username":"demo","password":"secret","base_path":"ClipBridgeServer"}'
```

Test the connection:

```bash
curl -X POST http://127.0.0.1:8787/api/admin/webdav/test \
  -H 'Authorization: Bearer new-admin-token'
```

Run one manual sync:

```bash
curl -X POST http://127.0.0.1:8787/api/admin/webdav/sync \
  -H 'Authorization: Bearer new-admin-token'
```

## Web UI Usage

Open `http://127.0.0.1:8787/` and:

1. paste either the admin token or a device token into the token box
2. press `Save`
3. use `Quick Clipboard` to upload text, links, images, and small files
4. use `Copy Latest` to write the current server text or link into the browser clipboard
5. use `Favorites` and `Clipboard History` to organize mixed item types
6. use `Runtime Settings`, `Limits`, `Retention Status`, and `Cleanup Policy` with the admin token when operating the server long-term
7. use `WebDAV Sync` to save credentials, test the connection, and run one manual sync
8. use `Pair Devices` and `Paired Devices` with the admin token when onboarding clients

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
- `internal/webdav/service.go`: persisted WebDAV settings, connection test, manual sync, remote manifest handling, and import/export logic
- `internal/store/sqlite.go`: SQLite connection lifecycle plus clipboard, favorites, categories, settings, and cleanup storage methods
- `internal/store/settings.go`: cleanup policy validation helpers
- `internal/store/sync_items.go`: import path for clipboard items pulled back from WebDAV
- `internal/store/webdav_settings.go`: persisted WebDAV settings and sync status storage helpers
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
- `internal/api/media_handler.go`: image, file, link, and download HTTP handlers
- `internal/api/cleanup_handler.go`: manual cleanup, storage status, and cleanup settings handlers
- `internal/api/settings_handler.go`: runtime settings and limits handlers
- `internal/api/webdav_handler.go`: WebDAV settings, test, sync, and status handlers
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
- `docs/deployment.md`: English deployment guide for binary, Docker, and service installs
- `docs/deployment.zh-CN.md`: 中文部署说明，方便直接给中文用户查看
- `docs/webdav-sync.md`: WebDAV sync preview scope, layout, and workflow
- `CHANGELOG.md`: phase-by-phase change record

## Notes

The Go module path is normalized to:

`github.com/cinmou/ClipBridgeServer`

All internal imports should use that exact casing so builds stay consistent
across macOS, Windows, Linux, and OpenWrt.
