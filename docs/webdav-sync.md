# WebDAV Sync Preview

Phase 11 adds a manual WebDAV sync preview so ClipBridgeServer can exchange
clipboard history and file payloads with storage the user already owns.

Current scope:

- configure WebDAV in the embedded Web UI or through `/api/settings/webdav`
- test the connection without restarting the server
- run one manual sync through `/api/admin/webdav/sync`
- push local clipboard history metadata to WebDAV
- push local image and file payloads to WebDAV
- pull remote clipboard history metadata back into SQLite
- pull remote image and file payloads back into the local `data/` directory
- record the latest test and sync status in the `settings` table

Current non-goals for this preview:

- background automatic WebDAV sync
- delete propagation
- full multi-writer conflict resolution
- encryption of WebDAV payloads

## Remote Layout

The server writes this structure under the configured WebDAV base path:

```text
ClipBridgeServer/
  manifest.json
  items/
    <sync-key>.json
  files/
    <sync-key>.bin
```

`sync-key` is a deterministic SHA-256 based identifier derived from the
clipboard item content and timestamps so repeated manual syncs can skip
duplicates without a schema migration.

## API Endpoints

- `GET /api/settings/webdav`
- `PATCH /api/settings/webdav`
- `POST /api/admin/webdav/test`
- `POST /api/admin/webdav/sync`
- `GET /api/admin/webdav/status`

All of them require the admin token.

## Web UI Flow

1. open `GET /`
2. paste the admin token
3. scroll to `WebDAV Sync`
4. fill `URL`, `Username`, `Password`, and optional `Base Path`
5. save settings
6. run `Test Connection`
7. run `Run Sync`

The panel shows:

- last connection test time and result
- last sync time and last successful sync time
- pushed and pulled item counts
- pushed and pulled file counts
- the latest error message or success summary

## Conflict Handling

This preview keeps conflict handling intentionally simple:

- if a local item already maps to the same `sync-key`, it is skipped
- if a remote item uses a `sync-key` the local server has never seen, it is imported
- the preview does not attempt to merge edits into the same clipboard record

That is enough for the first self-hosted sync preview while keeping the server
single-binary and easy to reason about.
