# ClipBridgeServer

A lightweight self-hosted clipboard sync server with an embedded Web UI.

一个轻量、自托管、带内置 Web 管理界面的剪贴板同步服务端。

## Status

- Current version: `v0.2.0-beta.1` `Cherwell`
- Project status: beta
- `Cherwell` focuses on the embedded Web UI, practical self-hosting, and a light security baseline.

## Features

- Embedded Web UI
- Clipboard history
- Text, link, image, and file support
- Favorites
- Devices and pairing codes
- Cleanup and retention policy
- WebDAV manual sync preview
- Lightweight self-hosted deployment
- Optional TLS support
- Simple token-based access
- English and Simplified Chinese Web UI

## Screenshots

Screenshots are not included in this repository yet.

## Architecture

`ClipBridgeServer` is the core server process. It exposes the HTTP API, stores clipboard history in SQLite, manages pairing and device tokens, and serves the embedded Web UI from the same binary.

The Web UI is embedded into the server binary and is meant for day-to-day management and lightweight browser-based clipboard usage.

Desktop and mobile clients are planned as separate future projects. Those clients should stay thin and use the HTTP API rather than reimplement server logic locally.

## Quick Start

1. Copy the example config:

```bash
cp configs/config.example.yaml config.yaml
```

2. Build the embedded Web UI:

```bash
cd web
npm run build
cd ..
```

3. Build the server:

```bash
go build -o clipbridge-server ./cmd/server
```

4. Run the server:

```bash
./clipbridge-server -config config.yaml
```

5. Open the browser:

```text
http://127.0.0.1:8787/
```

## Configuration

Use [`configs/config.example.yaml`](configs/config.example.yaml) as the starting point.

Important settings:

- `server.host` and `server.port`
- `storage.data_dir`
- `storage.database_path`
- `auth.token`
- `tls.enabled`
- `tls.cert_file`
- `tls.key_file`

Notes:

- The default bind address is `127.0.0.1`.
- If `auth.token` is empty, ClipBridgeServer will generate an admin token and store it under `data/secrets/admin_token`.
- Built-in TLS is optional. You can keep HTTP for trusted local or LAN use, or enable TLS if you already have certificates.

## Web UI

The embedded Web UI includes:

- Dashboard: latest clipboard item, quick upload, and status summary
- History: reverse chronological clipboard browsing with search and type filters
- Favorites: protected items excluded from automatic cleanup
- Devices: pairing codes and paired device management
- Settings: session token, language switcher, runtime status, security notes, cleanup controls
- WebDAV: configuration, connection testing, and manual sync preview

Language support:

- English
- Simplified Chinese (`zh-CN`)

Related docs:

- [`docs/webui.md`](docs/webui.md)
- [`docs/api.md`](docs/api.md)
- [`docs/config.md`](docs/config.md)
- [`docs/deployment.md`](docs/deployment.md)
- [`docs/security.md`](docs/security.md)
- [`docs/zh-CN/overview.md`](docs/zh-CN/overview.md)

## Security Model

ClipBridgeServer uses a practical self-hosted security model rather than a heavy enterprise access stack.

- Default bind address is `127.0.0.1`
- Admin token can be auto-generated if missing
- Pairing codes are short-lived and single-use
- Device tokens are generated during pairing
- Logs should not expose tokens or WebDAV passwords
- Plain HTTP is acceptable for trusted LAN use
- For remote or public access, use an HTTPS reverse proxy
- Built-in TLS is optional for advanced users

More detail:

- [`docs/security.md`](docs/security.md)

## WebDAV Sync

WebDAV sync is an optional preview feature.

- It is designed as a user-owned storage path
- It is not dependent on any ClipBridge cloud service
- Current workflow is manual test plus manual sync, not automatic background sync

More detail:

- [`docs/api.md`](docs/api.md)
- [`docs/webdav-sync.md`](docs/webdav-sync.md)

## Development

Build the frontend:

```bash
cd web
npm run build
cd ..
```

Run tests:

```bash
go test ./...
```

Build the server:

```bash
go build ./cmd/server
```

## Roadmap

Near-term priorities:

- polish the `Cherwell` Web UI
- improve documentation
- prepare future desktop and mobile clients that use the HTTP API

## License

This project uses the existing repository license: `GPL-3.0-only`.

See [`LICENSE`](LICENSE).

## 中文说明

### 这是什么

ClipBridgeServer 是一个自托管的剪贴板同步服务端，支持文本、链接、图片、文件历史记录，并带有内置 Web 管理界面。

### 当前能做什么

- 管理剪贴板历史记录
- 收藏重要内容
- 通过配对码添加设备
- 配置清理策略
- 预览 WebDAV 手动同步能力
- 在浏览器里使用英文或简体中文界面

### 为什么服务端和客户端分离

服务端负责存储、同步、设备管理和 Web UI。未来的桌面端和移动端建议保持轻量，通过 HTTP API 与服务端交互，这样服务端能力可以集中演进，客户端也更容易维护。

### 内网使用建议

- 默认监听 `127.0.0.1`
- 如果你只在可信内网中使用，HTTP 就已经足够实用
- 如果暴露到局域网，请确认网络环境可信

### 外网访问建议

如果需要远程或公网访问，建议放在 HTTPS 反向代理之后，例如 Caddy、Nginx、Traefik、Tailscale，或 NAS / 路由器自带的 HTTPS 网关。
