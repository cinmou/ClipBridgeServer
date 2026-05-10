# 部署说明

ClipBridgeServer 到第十阶段的目标，不再只是让开发者在 IDE 里
`go run`，而是让普通用户也能稳定安装、启动、常驻运行和更新。

项目的部署原则仍然保持不变：

- 一个服务端二进制
- 一个 `config.yaml`
- 一个 `data/` 目录

## 当前提供的部署资产

仓库里已经补好了这些内容：

- `scripts/build-release.sh`
  用来本地交叉编译多平台二进制，输出到 `dist/`
- `.github/workflows/release.yml`
  用来在 GitHub Actions 里自动测试、自动构建、自动上传 Release 产物
- `Dockerfile`
  用来构建容器镜像
- `docker-compose.example.yml`
  用来快速启动 Docker 版服务
- `deploy/systemd/clipbridge-server.service`
  Linux `systemd` 示例
- `deploy/launchd/com.cinmou.clipbridge-server.plist`
  macOS `launchd` 示例
- `deploy/windows/nssm-install.ps1`
  Windows + NSSM 常驻运行示例
- `deploy/openwrt/clipbridge-server.init`
  OpenWrt `init.d` 示例
- `deploy/homebrew/clipbridge-server.rb`
  Homebrew Tap 示例 Formula
- `deploy/scoop/clipbridge-server.json`
  Scoop Bucket 示例 Manifest

## 发布产物

当前发布流程会构建这些二进制：

- `clipbridge-server-linux-amd64`
- `clipbridge-server-linux-arm64`
- `clipbridge-server-darwin-amd64`
- `clipbridge-server-darwin-arm64`
- `clipbridge-server-windows-amd64.exe`

如果推送 `v*` 版本标签，GitHub Actions 会自动把这些产物上传到
GitHub Releases。

## 本地构建

只构建当前平台：

```bash
go build -o clipbridge-server ./cmd/server
```

构建全部发布二进制：

```bash
bash scripts/build-release.sh
```

构建完成后，产物会出现在 `dist/` 目录。

## 最简运行方式

推荐目录结构：

```text
clipbridge/
  clipbridge-server
  config.yaml
  data/
```

启动命令：

```bash
./clipbridge-server -config ./config.yaml
```

首次启动后，服务会自动创建：

- `data/clipbridge.db`
- `data/uploads/` 等运行目录

浏览器打开：

```text
http://127.0.0.1:8787/
```

就可以访问内置 Web UI。

## Docker 部署

直接构建并运行：

```bash
docker build -t clipbridge-server .
docker run --rm -p 8787:8787 \
  -v "$(pwd)/data:/app/data" \
  -v "$(pwd)/config.yaml:/app/config.yaml:ro" \
  clipbridge-server
```

或者使用仓库自带示例：

```bash
cp docker-compose.example.yml docker-compose.yml
docker compose up -d
```

这条路最适合 NAS、轻量服务器和家用 Linux 主机。

## Linux systemd

示例文件：

- `deploy/systemd/clipbridge-server.service`

典型安装步骤：

```bash
sudo useradd --system --home /opt/clipbridge --shell /usr/sbin/nologin clipbridge
sudo mkdir -p /opt/clipbridge/data
sudo cp clipbridge-server /opt/clipbridge/clipbridge-server
sudo cp config.yaml /opt/clipbridge/config.yaml
sudo cp deploy/systemd/clipbridge-server.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now clipbridge-server
```

这条路最适合 VPS、自建 Linux 服务器和长期后台运行。

## macOS launchd

示例文件：

- `deploy/launchd/com.cinmou.clipbridge-server.plist`

典型安装步骤：

```bash
sudo mkdir -p /usr/local/etc/clipbridge /usr/local/var/clipbridge/data /usr/local/var/log
sudo cp clipbridge-server /usr/local/bin/clipbridge-server
sudo cp config.yaml /usr/local/etc/clipbridge/config.yaml
cp deploy/launchd/com.cinmou.clipbridge-server.plist ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/com.cinmou.clipbridge-server.plist
```

如果以后做 Homebrew Tap，也可以直接配合 `brew services` 使用。

## Windows 常驻

示例脚本：

- `deploy/windows/nssm-install.ps1`

建议目录：

```text
C:\ClipBridge\
  clipbridge-server-windows-amd64.exe
  config.yaml
  data\
```

然后通过 NSSM 安装：

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\windows\nssm-install.ps1
```

这样就能把 ClipBridgeServer 当作长期后台服务运行。

## OpenWrt

示例文件：

- `deploy/openwrt/clipbridge-server.init`

建议路径：

- 二进制：`/usr/bin/clipbridge-server`
- 配置：`/etc/clipbridge/config.yaml`
- 工作目录：`/var/lib/clipbridge`

启用方式：

```sh
chmod +x /etc/init.d/clipbridge-server
/etc/init.d/clipbridge-server enable
/etc/init.d/clipbridge-server start
```

这条路主要是给路由器、软路由或轻量 OpenWrt 设备做参考。

## Homebrew 与 Scoop

这一步先提供可维护的模板，而不是直接在线发布：

- Homebrew 示例：`deploy/homebrew/clipbridge-server.rb`
- Scoop 示例：`deploy/scoop/clipbridge-server.json`

后续真正发版时，只需要把里面占位的下载地址和 SHA-256
替换成 GitHub Releases 的真实值即可。

## 验收建议

部署完成后，建议按这个顺序检查：

1. 服务是否能用 `-config config.yaml` 正常启动
2. 是否自动生成 `data/clipbridge.db`
3. `GET /api/health` 是否正常返回
4. 浏览器访问 `/` 是否能打开内置 Web UI
5. 带 token 的剪贴板接口是否仍然可用

## 现在项目到哪一步了

到第十阶段为止，ClipBridgeServer 已经具备这些能力：

- 文本、链接、图片、小文件剪贴板历史存储
- Bearer Token 鉴权
- 一次性配对码与长期 `device_token`
- 收藏夹与分类管理
- 清理策略与运行时设置
- 内置 Web UI
- 单二进制部署
- Docker / systemd / launchd / NSSM / OpenWrt 示例
- GitHub Actions 自动构建 Release 产物

这意味着项目已经从“开发者能跑”进入“普通用户可以部署体验”的阶段。
