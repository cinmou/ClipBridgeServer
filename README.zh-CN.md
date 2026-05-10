# ClipBridgeServer

一个轻量、自托管、跨平台的剪贴板服务端。它把 REST API、SQLite、
本地文件存储、设备配对、历史管理、内置 Web UI 和 WebDAV 同步预览
都放进同一个 Go 二进制里，部署目标始终保持为：

- 一个服务端二进制
- 一个 `config.yaml`
- 一个 `data/` 目录

语言入口：

- English: `README.md`
- 简体中文: `README.zh-CN.md`

## Beta 1

这份 README 现在作为 ClipBridgeServer 第一个正式 Beta 版入口文档。

这里说的 Beta 1，意思是：

- 服务端已经可以真实拿来试用
- 单二进制部署路径已经比较稳定
- 内置 Web UI 已经能承担日常管理工作
- 但仍然有一部分能力处于预览版或尚未交付状态

## 当前进度

当前已经推进到第 11 阶段：

- 第 8 阶段：图片、文件、链接类型支持
- 第 9 阶段：Settings API + Web Settings
- 第 10 阶段：Deployment
- 第 11 阶段：WebDAV Sync Preview

现在这个仓库已经能做这些事情：

- 作为局域网或自托管环境里的中心剪贴板服务
- 保存文本、链接、图片和小文件历史
- 通过配对码给客户端发长期 `device_token`
- 提供收藏夹和分类管理
- 提供自动清理和运行时限制配置
- 通过内置 Web UI 直接管理服务
- 手动把本地历史同步到用户自己的 WebDAV 存储

还没做的重点有：

- 桌面端或移动端自动后台监听剪贴板
- WebDAV 后台自动同步
- 更复杂的多端冲突合并
- 独立于 Bearer Token 的浏览器登录会话

## 后面要做什么

这个 Beta 之后，主线工作大致会放在这些方向：

- 桌面客户端后台监听本机剪贴板
- Windows / Linux 托盘与 macOS 菜单栏常驻
- WebDAV 从手动同步升级到后台自动同步
- 更清楚的同步日志和更稳的冲突处理
- 可选端到端加密，避免服务端看到明文剪贴板
- 更完善的浏览器和客户端登录流程
- 更适合长期自托管的导入导出与备份能力

## 快速开始

1. 复制示例配置：

```bash
cp configs/config.example.yaml config.yaml
```

2. 启动服务：

```bash
go run ./cmd/server -config config.yaml
```

3. 健康检查：

```bash
curl http://127.0.0.1:8787/api/health
```

预期返回：

```json
{"data":{"ok":true,"version":"0.1.0"}}
```

4. 打开内置 Web UI：

```text
http://127.0.0.1:8787/
```

服务启动后会自动创建：

- `data/clipbridge.db`

## 当前可以怎么用

### 1. 管理 Token

示例配置默认管理 Token 是：

```text
dev-token-123
```

除 `/api/health` 和 `/api/auth/pair` 外，受保护接口都需要：

```http
Authorization: Bearer <token>
```

### 2. 上传文本

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/text \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"text":"hello from ClipBridge"}'
```

Web UI 也复用同一个接口：

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/text \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"content":"from browser","source_device_id":"web-ui","source_device_name":"Web UI"}'
```

### 3. 上传链接和文件

上传链接：

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/link \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","source_device_id":"web-ui","source_device_name":"Web UI"}'
```

上传图片或文件：

```bash
curl -X POST http://127.0.0.1:8787/api/clipboard/file \
  -H 'Authorization: Bearer dev-token-123' \
  -F 'file=@./demo.png' \
  -F 'source_device_id=web-ui' \
  -F 'source_device_name=Web UI'
```

### 4. 查看最新和历史

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/latest

curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/history
```

按分类过滤：

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  'http://127.0.0.1:8787/api/clipboard/history?category=text'
```

### 5. 收藏和分类

收藏一条记录：

```bash
curl -X POST \
  -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/clipboard/items/1/favorite
```

查看收藏：

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/favorites
```

修改分类：

```bash
curl -X PATCH http://127.0.0.1:8787/api/clipboard/items/1/category \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"category":"work"}'
```

### 6. 设备配对

生成一次性配对码：

```bash
curl -X POST http://127.0.0.1:8787/api/auth/pairing-codes \
  -H 'Authorization: Bearer dev-token-123'
```

客户端换取长期 `device_token`：

```bash
curl -X POST http://127.0.0.1:8787/api/auth/pair \
  -H 'Content-Type: application/json' \
  -d '{"pairing_code":"ABCDEFGH","device_name":"My Laptop"}'
```

之后客户端就可以用返回的 `device_token` 调用剪贴板接口，不需要用户手动复制长管理 Token。

### 7. 清理策略和运行时设置

查看清理策略：

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/settings/cleanup
```

手动触发一次清理：

```bash
curl -X POST http://127.0.0.1:8787/api/admin/cleanup/run \
  -H 'Authorization: Bearer dev-token-123'
```

查看运行时设置：

```bash
curl -H 'Authorization: Bearer dev-token-123' \
  http://127.0.0.1:8787/api/settings
```

### 8. WebDAV 同步预览

保存 WebDAV 配置：

```bash
curl -X PATCH http://127.0.0.1:8787/api/settings/webdav \
  -H 'Authorization: Bearer dev-token-123' \
  -H 'Content-Type: application/json' \
  -d '{"enabled":true,"url":"https://dav.example.com/remote.php/dav/files/user","username":"demo","password":"secret","base_path":"ClipBridgeServer"}'
```

测试连接：

```bash
curl -X POST http://127.0.0.1:8787/api/admin/webdav/test \
  -H 'Authorization: Bearer dev-token-123'
```

手动同步：

```bash
curl -X POST http://127.0.0.1:8787/api/admin/webdav/sync \
  -H 'Authorization: Bearer dev-token-123'
```

这个阶段的 WebDAV 还是预览版，重点是：

- 用户可以填自己的 WebDAV
- 可以测试连接
- 可以手动推送和拉回历史
- 服务端不依赖我们的云服务

还没有做：

- 后台自动同步
- 删除同步
- 更复杂的冲突处理

## Web UI 使用方式

打开 `http://127.0.0.1:8787/` 之后：

1. 在顶部填入 admin token 或 device token
2. 用 `Quick Clipboard` 上传文本、链接、图片和文件
3. 用 `Copy Latest` 把服务端最新文本复制回浏览器剪贴板
4. 在 `Clipboard History` 和 `Favorites` 里管理历史
5. 用 `Runtime Settings`、`Limits`、`Cleanup Policy` 管理运行时参数
6. 用 `Pair Devices` 和 `Paired Devices` 管理设备
7. 用 `WebDAV Sync` 保存配置、测试连接和手动同步

注意：

- admin token 可以访问配对、设备管理、清理和设置接口
- device token 主要用于剪贴板、收藏夹和分类接口
- 浏览器写本地剪贴板必须由用户点击触发
- token 只保存在当前浏览器的本地存储里

## 部署

如果你只是想把它跑起来，最短路径是：

1. 从 GitHub Releases 下载对应平台二进制
2. 放到 `config.yaml` 同目录
3. 建一个空的 `data/` 目录
4. 运行 `./clipbridge-server -config ./config.yaml`

部署文档见：

- `docs/deployment.md`
- `docs/deployment.zh-CN.md`

## 我怎么构建版本

如果你只是想在当前机器构建一个可运行二进制：

```bash
go build -o clipbridge-server ./cmd/server
```

如果你想直接构建这一版 Beta 的多平台发布产物：

```bash
bash scripts/build-release.sh
```

构建完成后，产物会出现在 `dist/` 目录：

- `clipbridge-server-linux-amd64`
- `clipbridge-server-linux-arm64`
- `clipbridge-server-darwin-amd64`
- `clipbridge-server-darwin-arm64`
- `clipbridge-server-windows-amd64.exe`

推荐你本地按这个顺序做：

1. 先跑测试：

```bash
env GOCACHE=$(pwd)/.gocache go test ./...
```

2. 再构建发布产物：

```bash
bash scripts/build-release.sh
```

3. 打开 `dist/`，按目标平台取对应二进制。

如果你想走 GitHub Releases 的自动构建流程：

1. 提交代码
2. 打一个版本标签，比如 `v0.11.0-beta1`
3. 推送这个标签
4. 让 `.github/workflows/release.yml` 自动测试、自动构建、自动上传产物

## 更多文档

- `docs/api.md`
- `docs/config.md`
- `docs/architecture.md`
- `docs/roadmap.md`
- `docs/webdav-sync.md`
- `CHANGELOG.md`
