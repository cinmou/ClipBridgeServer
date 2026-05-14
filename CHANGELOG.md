# Changelog

## v0.2.0-beta.1 "Cherwell"

Second beta release focused on the embedded Web UI upgrade.

- redesigned the embedded Web UI into a full clipboard client and management console
- introduced the Web UI redesign initiative under the `WebDock` codename
- added a dashboard quick clipboard workflow for browser text upload, latest fetch, and browser clipboard copy
- expanded history and favorites with richer cards, detail views, search, filtering, file metadata, thumbnails, and item actions
- refreshed pairing, settings, cleanup, storage, and WebDAV management screens for consistent mobile-friendly use
- made embedded frontend routing refresh-safe for path-based navigation inside the single Go binary
- exposed clipboard item `size_bytes` in API responses for richer file and image metadata in the Web UI

## v0.1.0-beta.1 "Lea"

First public beta release of ClipBridgeServer.

## 0.0.11

- added persisted WebDAV settings through `GET /api/settings/webdav` and `PATCH /api/settings/webdav`
- added `POST /api/admin/webdav/test`, `POST /api/admin/webdav/sync`, and `GET /api/admin/webdav/status`
- added manual WebDAV sync preview for text, link, image, and file clipboard items
- added deterministic sync keys and import helpers for remote clipboard items
- updated the embedded Web UI with WebDAV settings, connection testing, and manual sync
- added `docs/webdav-sync.md`

## 0.0.10

- added `scripts/build-release.sh` for local multi-platform release builds
- added `.github/workflows/release.yml` for test, build, and tagged GitHub Release publishing
- added `Dockerfile` and `docker-compose.example.yml`
- added `systemd`, `launchd`, Windows NSSM, and OpenWrt service examples
- added Homebrew Tap and Scoop template manifests
- added `docs/deployment.md` and `docs/deployment.zh-CN.md`

## 0.0.9

- added `GET /api/settings` and `PATCH /api/settings`
- added `GET /api/settings/limits` and `PATCH /api/settings/limits`
- made admin token and request limits runtime-persisted settings
- updated the embedded Web UI with runtime settings and limits panels

## 0.0.8

- added link clipboard uploads through `POST /api/clipboard/link`
- added image and file uploads through `POST /api/clipboard/file`
- added file and image downloads through `GET /api/clipboard/items/{id}/file`
- stored file metadata in SQLite while streaming binary content to disk
- added SHA-256 hashing and MIME detection for uploaded files
- expanded clipboard history and latest endpoints to cover text, image, link, and file items
- updated the embedded Web UI to distinguish mixed clipboard item types

## 0.0.7

- added a Web UI quick clipboard panel that reuses the existing clipboard API
- added support for `content` and source metadata in clipboard uploads
- included source metadata in clipboard item responses
- added cleanup policy persistence in the `settings` table
- added a background cleaner worker with TTL, max item count, and max storage enforcement
- added favorite-aware cleanup rules
- added `GET /api/settings/cleanup` and `PATCH /api/settings/cleanup`
- added `POST /api/admin/cleanup/run`
- added `GET /api/admin/cleanup/status`
- added `GET /api/admin/storage/status`
- added Web UI panels for retention status and cleanup policy editing
- updated README and docs for phases 6.5 and 7

## 0.0.6

- added favorite state to clipboard records
- added favorite and unfavorite APIs for clipboard items
- added `GET /api/favorites`
- seeded built-in categories `text`, `image`, `link`, and `file`
- added category list and custom category creation APIs
- added per-item category reassignment and history filtering by category
- updated the embedded Web UI with favorites, category controls, and history filters
- expanded store and API tests to cover favorites and categories
- updated `README.md` and `docs/api.md` for the phase 6 workflow

## 0.0.5

- added an embedded Web UI served from `GET /`
- embedded `web/dist` into the Go binary
- added a browser console for health, latest clipboard, history, pairing codes, and devices
- added browser actions for delete, copy, and pairing code generation
- preserved the single-binary deployment model

## 0.0.4.5

- added one-time pairing codes with 5 minute expiry
- added long-lived device tokens for clipboard clients
- added `devices` and `pairing_codes` tables
- stored pairing codes and device tokens as hashes only
- added device list and revoke APIs
- added pairing and device authentication tests
- reserved a pairing URI for future QR-based onboarding

## 0.0.4

- added bearer token protection for all non-health API routes
- changed success responses to use a top-level `data` envelope
- changed error responses to use a top-level `error` envelope
- added request-size protection through `limits.max_request_bytes`
- added basic CORS support for local Web UI development
- added basic request logging with method, path, status code, and duration
- documented token usage, response envelopes, and error behavior

## 0.0.3

- added the text clipboard API for phase 3
- added `POST /api/clipboard/text`
- added `GET /api/clipboard/latest`
- added `GET /api/clipboard/history`
- added `GET /api/clipboard/items/{id}`
- added `DELETE /api/clipboard/items/{id}`
- added text size validation through `limits.min_text_bytes` and `limits.max_text_bytes`
- added SQLite-backed text clipboard CRUD methods
- added store and API tests for the new text flow
- documented the single-binary plus `config.yaml` plus `data/` deployment principle

## 0.0.2

- added the SQLite foundation layer for phase 2
- restored `storage.data_dir` and `storage.database_path` as active config keys
- initialized SQLite during server startup before the HTTP listener starts
- added migration tracking through `schema_migrations`
- added `migrations/001_init.sql` to create base tables
- verified `data/clipbridge.db` is created automatically
- kept `GET /api/health` unchanged after adding SQLite

## 0.0.1

- created the minimal Go server skeleton for stable startup
- normalized imports to match module path `github.com/cinmou/ClipBridgeServer`
- limited `cmd/server/main.go` to config loading, router creation, and HTTP startup
- reduced `internal/config` to configuration structure, defaults, and validation
- added `internal/api` route registration and `GET /api/health`
- verified `go run ./cmd/server -config config.yaml` and `curl http://127.0.0.1:8787/api/health`
