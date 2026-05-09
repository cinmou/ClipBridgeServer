# Changelog

## 0.7.0

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

## 0.6.0

- added favorite state to clipboard records
- added favorite and unfavorite APIs for clipboard items
- added `GET /api/favorites`
- seeded built-in categories `text`, `image`, `link`, and `file`
- added category list and custom category creation APIs
- added per-item category reassignment and history filtering by category
- updated the embedded Web UI with favorites, category controls, and history filters
- expanded store and API tests to cover favorites and categories
- updated `README.md` and `docs/api.md` for the phase 6 workflow

## 0.5.0

- added an embedded Web UI served from `GET /`
- embedded `web/dist` into the Go binary
- added a browser console for health, latest clipboard, history, pairing codes, and devices
- added browser actions for delete, copy, and pairing code generation
- preserved the single-binary deployment model

## 0.4.5

- added one-time pairing codes with 5 minute expiry
- added long-lived device tokens for clipboard clients
- added `devices` and `pairing_codes` tables
- stored pairing codes and device tokens as hashes only
- added device list and revoke APIs
- added pairing and device authentication tests
- reserved a pairing URI for future QR-based onboarding

## 0.4.0

- added bearer token protection for all non-health API routes
- changed success responses to use a top-level `data` envelope
- changed error responses to use a top-level `error` envelope
- added request-size protection through `limits.max_request_bytes`
- added basic CORS support for local Web UI development
- added basic request logging with method, path, status code, and duration
- documented token usage, response envelopes, and error behavior

## 0.3.0

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

## 0.2.0

- added the SQLite foundation layer for phase 2
- restored `storage.data_dir` and `storage.database_path` as active config keys
- initialized SQLite during server startup before the HTTP listener starts
- added migration tracking through `schema_migrations`
- added `migrations/001_init.sql` to create base tables
- verified `data/clipbridge.db` is created automatically
- kept `GET /api/health` unchanged after adding SQLite

## 0.1.0

- created the minimal Go server skeleton for stable startup
- normalized imports to match module path `github.com/cinmou/ClipBridgeServer`
- limited `cmd/server/main.go` to config loading, router creation, and HTTP startup
- reduced `internal/config` to configuration structure, defaults, and validation
- added `internal/api` route registration and `GET /api/health`
- verified `go run ./cmd/server -config config.yaml` and `curl http://127.0.0.1:8787/api/health`
