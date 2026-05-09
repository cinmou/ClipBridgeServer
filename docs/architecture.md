# Embedded Web UI Architecture

## Scope

This document describes the current server with SQLite startup, the text
clipboard API, the auth safety layer, the device pairing system, and the first
embedded browser-based management console.

At this stage the service is responsible for:

- loading `config.yaml`
- validating the server listen address
- validating the SQLite data directory and database path
- validating bearer token configuration
- validating text clipboard size limits
- validating whole request body size limits
- opening SQLite
- creating the initial schema automatically
- building an HTTP router
- exposing `GET /api/health`
- generating short-lived one-time pairing codes
- exchanging a pairing code for a long-lived device token
- listing paired devices
- revoking paired devices
- serving an embedded Web UI from `GET /`
- accepting text clipboard uploads
- reading the newest text clipboard item
- reading text clipboard history
- reading one text clipboard item by id
- soft-deleting one text clipboard item by id
- authenticating non-health API requests
- returning consistent JSON success and error envelopes
- allowing local development origins through CORS
- logging request method, path, status code, and duration
- storing pairing codes and device tokens as hashes only
- embedding built frontend assets into the Go binary
- starting the HTTP server

Everything else is postponed to later milestones.

## Request Flow

1. `cmd/server/main.go` parses `-config`
2. `internal/config.Load` reads YAML, applies defaults, and validates values
3. `internal/store.OpenSQLite` creates the data directory if needed
4. `internal/store.OpenSQLite` opens `data/clipbridge.db`
5. `internal/store.RunMigrations` creates the base tables from `migrations/001_init.sql`
6. `internal/store.RunMigrations` also creates `devices` and `pairing_codes`
7. `web/embedded.go` exposes the embedded frontend bundle
8. `internal/api.NewRouter` creates the HTTP router and middleware chain
9. Non-API requests fall through to the embedded Web UI handler
10. CORS middleware handles browser preflight requests for local development origins
11. Request-size middleware applies `limits.max_request_bytes`
12. Admin auth protects `/api/auth/pairing-codes` and `/api/auth/devices*`
13. Device auth allows `/api/clipboard/*` through either the admin token or a valid device token
14. `POST /api/auth/pairing-codes` creates a 5 minute one-time code
15. `POST /api/auth/pair` exchanges one code for one device token and invalidates the code immediately
16. `POST /api/clipboard/text` validates the request body and delegates writes to `internal/store`
17. `GET` and `DELETE` clipboard routes delegate reads and soft deletes to `internal/store`
18. Response helpers wrap successes in `data` and failures in `error`
19. Logging middleware prints method, path, status code, and duration
20. `http.ListenAndServe` starts the server

## Package Responsibilities

### `cmd/server`

Owns process startup only.

- parse flags
- load config
- initialize SQLite
- create router
- start the HTTP listener

It should not contain business logic.

### `internal/config`

Owns configuration concerns only.

- config structures
- default values
- YAML loading
- validation

It should not create routers, stores, or background jobs.

### `internal/auth`

Owns authentication concerns only.

- parse bearer tokens
- generate pairing codes
- generate device tokens
- hash secrets before they are stored

It should not know about SQLite or clipboard storage.

### `internal/store`

Owns persistence bootstrapping only.

- create the local database directory
- open the SQLite connection
- ping the database
- create the schema migrations table
- apply the initial schema
- create one-time pairing codes
- create and authenticate device tokens
- list and revoke devices
- create, list, read, and soft-delete text clipboard rows

It should not contain clipboard business endpoints yet.

### `internal/api`

Owns HTTP route registration and handlers only.

- register API endpoints
- implement health check response
- validate HTTP requests
- map store results to JSON responses
- coordinate pairing and device management flows
- route non-API traffic to the embedded UI bundle

It should not load config files directly or manage database setup.

## Files In This Stage

- `cmd/server/main.go`: boot sequence
- `internal/config/config.go`: config loading and validation
- `internal/config/defaults.go`: defaults for auth, server, storage, and limits
- `internal/auth/middleware.go`: bearer token middleware
- `web/embedded.go`: frontend asset embedding and serving
- `internal/store/device_pairing.go`: pairing and device persistence
- `internal/store/sqlite.go`: SQLite initialization and cleanup
- `internal/store/migrations.go`: migration bookkeeping and execution
- `internal/store/sqlite_test.go`: verifies database and tables are created
- `internal/store/device_pairing_test.go`: verifies pairing code lifecycle, hash storage, and revocation
- `internal/store/text_items_test.go`: verifies text clipboard CRUD in SQLite
- `migrations/001_init.sql`: creates `clipboard_items`, `categories`, `clipboard_item_categories`, and `settings`
- `migrations/002_device_pairing.sql`: creates `devices` and `pairing_codes`
- `internal/api/router.go`: route wiring
- `internal/api/health_handler.go`: `GET /api/health`
- `internal/api/auth_handler.go`: pairing and device management handlers
- `internal/api/clipboard_handler.go`: text clipboard handlers
- `internal/api/response.go`: shared JSON response helpers
- `internal/api/clipboard_handler_test.go`: HTTP integration-style tests
- `web/dist/*`: frontend build output embedded into the binary

## Acceptance Standard

The current stage is considered complete when these flows succeed:

```bash
go run ./cmd/server -config config.yaml
curl http://127.0.0.1:8787/api/health
curl -X POST http://127.0.0.1:8787/api/auth/pairing-codes -H 'Authorization: Bearer dev-token-123'
curl -X POST http://127.0.0.1:8787/api/auth/pair -H 'Content-Type: application/json' -d '{"pairing_code":"ABCDEFGH","device_name":"Laptop"}'
curl -X POST http://127.0.0.1:8787/api/clipboard/text -H 'Authorization: Bearer <device_token>' -H 'Content-Type: application/json' -d '{"text":"hello"}'
curl -H 'Authorization: Bearer <device_token>' http://127.0.0.1:8787/api/clipboard/latest
curl -H 'Authorization: Bearer dev-token-123' http://127.0.0.1:8787/api/auth/devices
curl http://127.0.0.1:8787/
```

And the following are true:

- `data/clipbridge.db` is created automatically
- base tables exist in SQLite
- pairing codes expire after 5 minutes and become invalid after one successful exchange
- device tokens are stored as hashes and can authenticate clipboard requests
- revoked devices lose clipboard access immediately
- the embedded Web UI is reachable from `/` without a separate frontend deployment
- text clipboard items are written into SQLite and can be queried back
- missing token returns `401`
- oversized request body returns `413`
- the health HTTP response body is:

```json
{"data":{"ok":true,"version":"0.1.0"}}
```
