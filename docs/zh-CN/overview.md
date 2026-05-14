# ClipBridgeServer 概览

## ClipBridgeServer 是什么

ClipBridgeServer 是一个轻量、自托管的剪贴板同步服务端，支持文本、链接、图片、文件历史记录，并带有内置 Web 管理界面。

## 当前版本能做什么

当前 `v0.2.0-beta.1 Cherwell` 版本已经支持：

- 内置 Web UI
- 剪贴板历史记录
- 文本、链接、图片、文件支持
- 收藏
- 配对码与设备管理
- 清理与保留策略
- WebDAV 手动同步预览
- 英文与简体中文界面

## 如何启动

```bash
cp configs/config.example.yaml config.yaml
go build -o clipbridge-server ./cmd/server
./clipbridge-server -config config.yaml
```

## 如何进入 Web UI

默认本机访问地址：

```text
http://127.0.0.1:8787/
```

## 如何使用访问令牌

- 管理接口使用 admin token
- 剪贴板客户端可通过配对流程获得 device token
- 如果 `auth.token` 为空，服务端会自动生成 admin token
- 自动生成的 token 会保存在 `data/secrets/admin_token`

## 内网使用建议

- 默认监听 `127.0.0.1`
- 如果只在本机或可信内网中使用，HTTP 是可以接受的
- 如果改成 `0.0.0.0`，请确认网络环境可信

## 外网访问建议

如果需要远程或公网访问，建议放在 HTTPS 反向代理之后，例如：

- Caddy
- Nginx
- Traefik
- Tailscale
- NAS / 路由器自带 HTTPS 网关

## WebDAV 同步说明

当前 WebDAV 是预览能力，偏向手动同步：

- 保存配置
- 测试连接
- 手动触发同步

它不是当前版本中的自动后台同步系统。

## 服务端和未来客户端的关系

ClipBridgeServer 是核心服务端，负责：

- 数据存储
- 剪贴板历史
- 配对与设备管理
- 内置 Web UI

未来桌面端和移动端建议保持轻量，通过 HTTP API 与服务端交互，而不是把服务端逻辑分散到各个客户端里。
