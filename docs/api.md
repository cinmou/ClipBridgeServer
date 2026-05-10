# API

## Authentication Model

There are now two credential types:

- admin token: configured in `config.yaml`, used for management APIs
- device token: returned by the pairing flow, used by clients for clipboard APIs

Rules:

- `/api/health` does not require a token
- `/api/auth/pair` does not require a token
- `/api/auth/pairing-codes`, `/api/auth/devices*`, `/api/admin/*`, and `/api/settings/*` require the admin token
- `/api/clipboard/*`, `/api/favorites`, and `/api/categories` accept either the admin token or a valid device token

Bearer format:

```http
Authorization: Bearer <token>
```

## Response Format

Successful responses use:

```json
{
  "data": {}
}
```

Failed responses use:

```json
{
  "error": {
    "message": "..."
  }
}
```

## Embedded Web UI

`GET /` serves the built Web UI from inside the Go binary.

Current phase 11 behavior:

- the browser saves the chosen admin token or device token in local storage
- the UI reuses that token for later API calls
- the UI can upload browser text, links, images, and files through existing clipboard APIs
- the UI can show latest clipboard, mixed-type history, favorites, pairing codes, devices, cleanup status, runtime settings, and WebDAV sync status
- the UI can manually copy latest server text into the browser clipboard
- the UI can save WebDAV settings, test the connection, and trigger one manual sync

## Health Check

### `GET /api/health`

Response:

```json
{
  "data": {
    "ok": true,
    "version": "0.1.0"
  }
}
```

## Device Pairing

### `POST /api/auth/pairing-codes`

Requires the admin token.

Creates one one-time pairing code that expires after 5 minutes.

### `POST /api/auth/pair`

Does not require an existing token.

Request:

```json
{
  "pairing_code": "ABCDEFGH",
  "device_name": "My Laptop"
}
```

### `GET /api/auth/devices`

Requires the admin token.

Returns all paired devices, including revoked ones.

### `DELETE /api/auth/devices/{id}`

Requires the admin token.

Revokes one device immediately.

## Clipboard Items

### Item Shape

Clipboard item responses now include favorite, category, and source metadata:

```json
{
  "id": 1,
  "type": "text",
  "text": "hello from ClipBridge",
  "is_favorite": false,
  "category": "text",
  "source_device_id": "web-ui",
  "source_device_name": "Web UI",
  "created_at": "2026-05-09T00:00:00Z",
  "updated_at": "2026-05-09T00:00:00Z"
}
```

### `POST /api/clipboard/text`

Stores one text clipboard item.

Classic request body:

```json
{
  "text": "hello from ClipBridge"
}
```

Web UI compatible request body, using the same API:

```json
{
  "content": "typed in browser",
  "source_device_id": "web-ui",
  "source_device_name": "Web UI"
}
```

### `POST /api/clipboard/link`

Stores one link clipboard item.

Request:

```json
{
  "url": "https://example.com",
  "source_device_id": "web-ui",
  "source_device_name": "Web UI"
}
```

### `POST /api/clipboard/file`

Stores one image or generic file clipboard item through `multipart/form-data`.

Form fields:

- `file`: binary file payload
- `source_device_id`: optional source identifier
- `source_device_name`: optional source name

Possible errors:

- `401` when the token is missing, invalid, or revoked
- `400` when the body is empty, invalid JSON, or text size is outside limits
- `413` when the request body exceeds `limits.max_request_bytes`

### `GET /api/clipboard/latest`

Returns the newest non-deleted clipboard item across all supported types.

### `GET /api/clipboard/history`

Returns clipboard items ordered from newest to oldest.

Optional filter:

```text
GET /api/clipboard/history?category=text
```

### `GET /api/clipboard/items/{id}`

Returns one non-deleted text clipboard item by id.

### `DELETE /api/clipboard/items/{id}`

Soft-deletes one text clipboard item by id.

### `GET /api/clipboard/items/{id}/file`

Streams one stored image or file payload back to the client.

Notes:

- images are served inline for preview-capable clients
- generic files are served as downloadable attachments

## Favorites

### `POST /api/clipboard/items/{id}/favorite`

Marks one clipboard item as a favorite.

### `DELETE /api/clipboard/items/{id}/favorite`

Removes one clipboard item from favorites.

### `GET /api/favorites`

Returns all non-deleted favorite clipboard items ordered from newest to oldest.

Notes:

- favorites are excluded from automatic cleanup

## Categories

Built-in categories seeded by migrations:

- `text`
- `image`
- `link`
- `file`

### `GET /api/categories`

Returns all built-in and custom categories.

### `POST /api/categories`

Creates one custom category.

Request:

```json
{
  "name": "work"
}
```

### `PATCH /api/clipboard/items/{id}/category`

Reassigns one clipboard item to one category.

Request:

```json
{
  "category": "work"
}
```

## Cleanup And Storage

### `GET /api/settings/cleanup`

Requires the admin token.

Returns the persisted cleanup policy.

Response:

```json
{
  "data": {
    "ttl_hours": 168,
    "max_items": 1000,
    "max_total_size_mb": 2048,
    "interval_minutes": 30,
    "enabled": true
  }
}
```

### `PATCH /api/settings/cleanup`

Requires the admin token.

Updates the cleanup policy and makes it effective without restarting the
service.

Request:

```json
{
  "ttl_hours": 168,
  "max_items": 1000,
  "max_total_size_mb": 2048,
  "interval_minutes": 30,
  "enabled": true
}
```

### `GET /api/admin/cleanup/status`

Requires the admin token.

Returns the latest cleanup run summary.

### `POST /api/admin/cleanup/run`

Requires the admin token.

Triggers one cleanup run immediately.

### `GET /api/admin/storage/status`

Requires the admin token.

Returns the current active storage summary:

- `history_count`
- `favorite_count`
- `total_bytes`
- `file_bytes`

## Runtime Settings

### `GET /api/settings`

Requires the admin token.

Returns the combined runtime settings and startup-only settings.

### `PATCH /api/settings`

Requires the admin token.

Current implementation supports changing:

- `admin_token`

### `GET /api/settings/webdav`

Requires the admin token.

Returns the persisted WebDAV sync preview settings.

### `PATCH /api/settings/webdav`

Requires the admin token.

Updates the persisted WebDAV sync preview settings.

Request:

```json
{
  "enabled": true,
  "url": "https://dav.example.com/remote.php/dav/files/user",
  "username": "demo",
  "password": "secret",
  "base_path": "ClipBridgeServer"
}
```

### `GET /api/settings/limits`

Requires the admin token.

Returns the persisted runtime size limits.

### `PATCH /api/settings/limits`

Requires the admin token.

Updates:

- text min and max bytes
- image min and max bytes
- file min and max bytes
- link min and max bytes
- max request bytes

## WebDAV Sync Preview

### `POST /api/admin/webdav/test`

Requires the admin token.

Tests the current persisted WebDAV settings and records the latest test result.

### `POST /api/admin/webdav/sync`

Requires the admin token.

Runs one manual WebDAV sync pass:

- pushes local clipboard metadata
- pushes local image and file payloads
- pulls remote items not seen locally yet
- stores the latest sync summary in the settings table

### `GET /api/admin/webdav/status`

Requires the admin token.

Returns the latest connection test and sync summary.

## Safety Limits

Configured request safety that already applies in this phase:

- `limits.max_request_bytes`: max HTTP request body size
- `limits.min_text_bytes`: minimum text payload size
- `limits.max_text_bytes`: maximum text payload size

When the request body is too large, the API returns:

```json
{
  "error": {
    "message": "request body is too large"
  }
}
```
