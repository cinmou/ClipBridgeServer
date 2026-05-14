# API Overview

ClipBridgeServer exposes a small HTTP API used by the embedded Web UI and
future clients.

This is a concise overview, not a full OpenAPI document.

## Authentication Model

Current credential types:

- admin token
- device token

Bearer format:

```http
Authorization: Bearer <token>
```

Rules:

- public endpoints do not require a token
- admin management endpoints require the admin token
- clipboard endpoints allow either the admin token or a valid device token

## Response Format

Successful responses:

```json
{
  "data": {}
}
```

Failed responses:

```json
{
  "error": {
    "message": "..."
  }
}
```

## Public Endpoints

- `GET /api/health`
- `POST /api/auth/pair`

Notes:

- `GET /api/health` reports service health and version
- `POST /api/auth/pair` exchanges a valid pairing code for a device token

## Admin Token Required

### Pairing And Devices

- `POST /api/auth/pairing-codes`
- `GET /api/auth/devices`
- `DELETE /api/auth/devices/{id}`

### Settings

- `GET /api/settings`
- `PATCH /api/settings`
- `GET /api/settings/limits`
- `PATCH /api/settings/limits`
- `GET /api/settings/cleanup`
- `PATCH /api/settings/cleanup`
- `GET /api/settings/webdav`
- `PATCH /api/settings/webdav`

### Cleanup And Storage

- `POST /api/admin/cleanup/run`
- `GET /api/admin/cleanup/status`
- `GET /api/admin/storage/status`

### WebDAV Preview

- `POST /api/admin/webdav/test`
- `POST /api/admin/webdav/sync`
- `GET /api/admin/webdav/status`

## Device Token Allowed

These endpoints accept either:

- admin token
- valid device token

### Clipboard

- `POST /api/clipboard/text`
- `POST /api/clipboard/link`
- `POST /api/clipboard/file`
- `GET /api/clipboard/latest`
- `GET /api/clipboard/history`
- `GET /api/clipboard/items/{id}`
- `DELETE /api/clipboard/items/{id}`
- `GET /api/clipboard/items/{id}/file`
- `POST /api/clipboard/items/{id}/favorite`
- `DELETE /api/clipboard/items/{id}/favorite`
- `PATCH /api/clipboard/items/{id}/category`

### Favorites And Categories

- `GET /api/favorites`
- `GET /api/categories`
- `POST /api/categories`

## Clipboard Notes

Supported clipboard item types:

- text
- link
- image
- file

History is returned in reverse chronological order.

Favorites are excluded from automatic cleanup.

## Pairing Notes

- pairing codes expire after 5 minutes
- pairing codes are single-use
- device tokens are generated during successful pairing

## WebDAV Notes

WebDAV support is currently:

- preview/manual sync oriented
- configured through settings endpoints
- not a background always-on cloud sync system
