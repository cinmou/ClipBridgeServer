# Roadmap

## Phase 1: Server Foundation

Status: done

Delivered in this phase:

- normalized Go module and import path casing
- minimal `cmd/server` startup flow
- isolated `internal/config` package
- isolated `internal/api` package
- `GET /api/health`
- verified local boot path

Not included in this phase:

- SQLite
- auth middleware
- clipboard CRUD
- favorites
- categories
- TTL cleanup
- WebDAV

## Phase 2: Local Persistence

Status: done

- add SQLite initialization
- create schema and migrations
- wire store lifecycle into server startup

Delivered in this phase:

- added SQLite driver dependency
- created `internal/store` startup layer
- created `migrations/001_init.sql`
- auto-created `data/clipbridge.db` on startup
- auto-created `clipboard_items`, `categories`, `clipboard_item_categories`, `settings`, and `schema_migrations`
- kept `GET /api/health` working after SQLite startup

Not included in this phase:

- token middleware
- clipboard upload and download
- favorites behavior
- category business APIs
- TTL cleanup worker
- WebDAV sync

## Phase 3: Text Clipboard API

Status: done

Delivered in this phase:

- added text clipboard CRUD endpoints
- added request and response models for text items
- added text size validation through config
- kept deployment as one binary plus `config.yaml` plus `data/`
- added store tests and API tests

Not included in this phase:

- token middleware
- image clipboard APIs
- file clipboard APIs
- favorites
- categories
- WebDAV sync

## Phase 4: Auth And API Safety

Status: done

Delivered in this phase:

- added bearer token protection for every `/api/*` route except `/api/health`
- unified success responses under `data`
- unified failure responses under `error`
- added global request-size protection through `limits.max_request_bytes`
- added local-development CORS support
- added request logging
- documented the auth and error contract in `docs/api.md`

Not included in this phase:

- image clipboard APIs
- file clipboard APIs
- favorites
- categories
- TTL cleanup
- WebDAV sync

## Phase 4.5: Device Pairing

Status: done

Delivered in this phase:

- added one-time pairing codes with a 5 minute lifetime
- added device tokens for long-lived client authentication
- added `devices` and `pairing_codes` tables
- stored pairing codes and device tokens as hashes only
- added device list and revoke APIs
- reserved a pairing URI field for future QR onboarding

Not included in this phase:

- image clipboard APIs
- file clipboard APIs
- favorites
- categories
- TTL cleanup
- WebDAV sync

## Phase 5: Embedded Web UI MVP

Status: done

Delivered in this phase:

- added an embedded Web UI served from `GET /`
- embedded built frontend assets into the Go binary with `go:embed`
- showed health, latest clipboard, history, pairing codes, and paired devices in the browser
- added Web UI actions for delete, copy, and pairing code generation
- kept deployment as one server binary plus `config.yaml` plus `data/`

Not included in this phase:

- favorites UI
- category UI
- settings UI
- search and filtering
- rich authentication sessions

## Phase 6: Retention And Organization

Status: done

Delivered in this phase:

- added favorite state on clipboard records
- added favorite and unfavorite APIs
- added `GET /api/favorites`
- seeded built-in categories `text`, `image`, `link`, and `file`
- added category list and custom category creation APIs
- added item category reassignment
- added history filtering by category
- expanded the embedded Web UI with favorites, category filters, and category editing
- documented current progress and concrete usage in `README.md`

Not included in this phase:

- TTL cleanup worker behavior
- image clipboard APIs
- file clipboard APIs
- WebDAV sync

## Phase 6.5: Web UI Quick Clipboard Panel

Status: done

Delivered in this phase:

- added a browser-side quick clipboard panel
- reused `POST /api/clipboard/text` instead of creating a Web UI-only API
- added support for `content` plus source metadata in clipboard uploads
- allowed the Web UI to copy the latest server text into the browser clipboard
- surfaced upload source metadata in clipboard responses

## Phase 7: Cleaner And Retention Policy

Status: done

Delivered in this phase:

- added a background cleaner worker
- added persisted cleanup policy in the `settings` table
- added TTL cleanup for non-favorite items
- added max history count cleanup for oldest non-favorites
- added max storage size cleanup with oldest non-favorite file priority
- added manual cleanup, cleanup status, storage status, and cleanup settings APIs
- added embedded Web UI panels for quick clipboard, retention status, and cleanup policy editing
- kept cleanup settings live without requiring service restart

Not included in this phase:

- image clipboard APIs
- file clipboard APIs
- WebDAV sync

## Phase 8: Media Clipboard APIs

Planned:

- image clipboard APIs
- file clipboard APIs

## Phase 9: Remote Sync

Planned:

- WebDAV-based synchronization
