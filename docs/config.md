# Configuration

ClipBridgeServer reads runtime startup config from `config.yaml`.

Use [`configs/config.example.yaml`](../configs/config.example.yaml) as the
starting point.

## Main Fields

### `server.host`

Listen host for the embedded HTTP server.

Common values:

- `127.0.0.1`
- `0.0.0.0`
- `localhost`

Default behavior is local-only binding through `127.0.0.1`.

### `server.port`

Listen port for the server.

Typical value:

- `8787`

### `storage.data_dir`

Directory for runtime data such as:

- uploaded files
- generated secrets
- local sync/runtime state

### `storage.database_path`

Path to the SQLite database file.

Typical value:

- `./data/clipbridge.db`

### `auth.token`

Admin token for management access.

This field can be left empty.

If it is empty, ClipBridgeServer generates an admin token automatically and
stores it under:

```text
data/secrets/admin_token
```

### `tls.enabled`

Enables built-in TLS startup.

If `false`, the server uses normal HTTP startup.

### `tls.cert_file`

Certificate file path used when `tls.enabled` is `true`.

### `tls.key_file`

Private key file path used when `tls.enabled` is `true`.

## Cleanup And Retention

These startup values seed the runtime cleanup policy:

### `storage.ttl_hours`

Default retention TTL for non-favorite items.

### `storage.max_items`

Default maximum active clipboard history count.

### `storage.max_total_size_mb`

Default maximum active storage budget.

### `cleaner.enabled`

Whether the background cleanup worker should run.

### `cleaner.interval_minutes`

How often the cleanup worker wakes up and reloads policy.

## Request Limits

Current request and payload sizing fields:

- `limits.min_text_bytes`
- `limits.max_text_bytes`
- `limits.min_image_bytes`
- `limits.max_image_bytes`
- `limits.min_file_bytes`
- `limits.max_file_bytes`
- `limits.min_link_bytes`
- `limits.max_link_bytes`
- `limits.max_request_bytes`

## WebDAV Settings

WebDAV settings are runtime-managed rather than startup-managed.

They are updated through the Web UI or settings API and include:

- enabled flag
- server URL
- username
- password
- base path

These are not the main startup fields in `config.yaml`, but they are part of
the current server behavior.

## Example

```yaml
auth:
  token: ""

server:
  host: "127.0.0.1"
  port: 8787

tls:
  enabled: false
  cert_file: ""
  key_file: ""

storage:
  data_dir: "./data"
  database_path: "./data/clipbridge.db"
  ttl_hours: 168
  max_items: 1000
  max_total_size_mb: 2048

limits:
  min_text_bytes: 1
  max_text_bytes: 262144
  min_image_bytes: 1
  max_image_bytes: 10485760
  min_file_bytes: 1
  max_file_bytes: 52428800
  min_link_bytes: 1
  max_link_bytes: 8192
  max_request_bytes: 1048576

cleaner:
  enabled: true
  interval_minutes: 30
```
