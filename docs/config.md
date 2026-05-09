# Configuration

## Active Keys In The Current Implementation

The current code actively reads these keys from `config.yaml`:

### `auth.token`

Bearer token required by every protected API route.

### `server.host`

HTTP listen host.

Examples:

- `127.0.0.1`
- `0.0.0.0`
- `localhost`

### `server.port`

HTTP listen port.

Default:

- `8787`

### `storage.data_dir`

Directory used for local runtime data.

Default:

- `./data`

### `storage.database_path`

SQLite database file path.

Default:

- `./data/clipbridge.db`

### `storage.ttl_hours`

Default TTL used to seed the persisted cleanup policy.

### `storage.max_items`

Default maximum active history count used to seed the persisted cleanup policy.

### `storage.max_total_size_mb`

Default maximum active storage budget used to seed the persisted cleanup policy.

### `limits.min_text_bytes`

Minimum allowed text clipboard payload size in bytes.

### `limits.max_text_bytes`

Maximum allowed text clipboard payload size in bytes.

### `limits.max_request_bytes`

Maximum allowed HTTP request body size in bytes.

### `cleaner.enabled`

Whether the background cleanup worker should run by default.

### `cleaner.interval_minutes`

Default background cleanup interval in minutes.

## Current Validation Rules

- `auth.token` must not be empty
- `server.host` must be a valid IP address or `localhost`
- `server.port` must be between `1` and `65535`
- `storage.data_dir` must not be empty
- `storage.database_path` must not be empty
- `storage.ttl_hours` must be greater than `0`
- `storage.max_items` must be greater than `0`
- `storage.max_total_size_mb` must be greater than `0`
- `limits.min_text_bytes` must be greater than or equal to `0`
- `limits.max_text_bytes` must be greater than `0`
- `limits.max_request_bytes` must be greater than `0`
- `limits.min_text_bytes` must not be greater than `limits.max_text_bytes`
- `cleaner.interval_minutes` must be greater than `0`

## Example

```yaml
auth:
  token: "dev-token-123"

server:
  host: "0.0.0.0"
  port: 8787

storage:
  data_dir: "./data"
  database_path: "./data/clipbridge.db"
  ttl_hours: 168
  max_items: 1000
  max_total_size_mb: 2048

limits:
  min_text_bytes: 1
  max_text_bytes: 262144
  max_request_bytes: 1048576

cleaner:
  enabled: true
  interval_minutes: 30
```
