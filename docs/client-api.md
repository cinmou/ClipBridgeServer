# Client API

This document defines the minimal HTTP API surface future ClipBridge clients should use. Clients should authenticate with a `device_token` obtained through pairing, not with the admin token.

For the full server API overview, see [api.md](api.md).

## Authentication Model

Client pairing flow:

1. A user opens the Web UI and generates a pairing code.
2. The client sends the pairing code to `POST /api/auth/pair`.
3. The server returns a `device_token` and device metadata.
4. The client stores the device credentials locally.
5. Future client requests send:

```http
Authorization: Bearer <device_token>
```

Clients should store:

- `server_url`
- `device_id`
- `device_token`
- `device_name`
- `last_sync_at`

Clients should not call admin APIs such as settings, cleanup, WebDAV configuration, device revocation, or pairing-code generation.

## Public Endpoints

### `GET /api/health`

Use this to check whether the server is reachable before pairing or syncing.

Expected use:

- Validate `server_url`.
- Show connection status.
- Detect server version or basic health details if returned.

### `POST /api/auth/pair`

Exchange a short-lived pairing code for a device token.

Example request:

```json
{
  "pairing_code": "ABCD1234",
  "device_name": "Alice MacBook"
}
```

Expected client behavior:

- Send the pairing code only over trusted local networks or HTTPS.
- Store the returned `device_token` in platform secure storage when available.
- Treat failed pairing as recoverable and ask the user to generate a new code.

## Clipboard Endpoints

The following endpoints should be called with:

```http
Authorization: Bearer <device_token>
```

### `POST /api/clipboard/text`

Upload a text clipboard item.

Example request:

```json
{
  "content": "Hello from ClipBridge",
  "source_device_name": "Alice MacBook"
}
```

Clients may include local device metadata when supported by the server. They should keep payloads small and explicit in the MVP.

### `POST /api/clipboard/link`

Upload a link clipboard item, if this endpoint is available in the target server version.

Example request:

```json
{
  "url": "https://example.com",
  "title": "Example"
}
```

If a client cannot distinguish links from text reliably, it may upload the value as text instead.

### `POST /api/clipboard/file`

Upload an image or file clipboard item, if this endpoint is available in the target server version.

Expected format:

- Multipart form upload.
- File field should contain the file or image bytes.
- Optional source metadata may be included when supported.

Clients should avoid automatic large uploads in the MVP. Prefer an explicit user action before uploading files.

### `GET /api/clipboard/latest`

Fetch the latest clipboard item.

Expected use:

- Populate the home screen.
- Enable “download latest” behavior.
- Compare with `last_seen_item_id` to detect changes.

### `GET /api/clipboard/history`

Fetch recent clipboard history.

Expected use:

- Show a short recent-history list.
- Let users select an older item to copy, open, or download.

Clients should request only the amount of history needed for a lightweight native UI.

### `GET /api/clipboard/items/{id}`

Fetch details for one clipboard item.

Expected use:

- Open a detail view.
- Resolve metadata before download.
- Display content preview, MIME type, size, and source when available.

### `GET /api/clipboard/items/{id}/file`

Download the file payload for an item, if this endpoint is available and the item has file content.

Expected use:

- Download image/file items.
- Open or share downloaded content using native platform behavior.

## Admin APIs Are Out Of Scope

Clients should not call these API groups:

- Pairing-code generation.
- Device revocation.
- Server settings.
- Cleanup and storage administration.
- WebDAV configuration.
- WebDAV manual sync controls.

Those flows belong in the embedded Web UI.

## Compatibility Guidance

- Start with `GET /api/health` before pairing or sync.
- Treat optional endpoints as feature-detected behavior.
- Keep error messages human-readable and action-oriented.
- Do not assume future realtime sync exists.
- Do not require admin-token access for any client MVP flow.

