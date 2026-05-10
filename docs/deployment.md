# Deployment Guide

ClipBridgeServer is designed to stay easy to ship:

- one server binary
- one `config.yaml`
- one `data/` directory

This phase adds the packaging and service examples needed to move from
`go run` into normal user deployment.

## Release Artifacts

The release workflow builds these binaries:

- `clipbridge-server-linux-amd64`
- `clipbridge-server-linux-arm64`
- `clipbridge-server-darwin-amd64`
- `clipbridge-server-darwin-arm64`
- `clipbridge-server-windows-amd64.exe`

The GitHub Actions workflow lives in
`/.github/workflows/release.yml`. It runs tests, cross-builds the server, uploads
workflow artifacts, and publishes assets when a `v*` tag is pushed.

## Local Build

Build one local binary:

```bash
go build -o clipbridge-server ./cmd/server
```

Build all release binaries into `dist/`:

```bash
bash scripts/build-release.sh
```

## Minimal Runtime Layout

The recommended layout is:

```text
clipbridge/
  clipbridge-server
  config.yaml
  data/
```

Start it with:

```bash
./clipbridge-server -config ./config.yaml
```

After first startup the SQLite database and uploaded file storage stay under
`data/`.

## Docker

Build and run with Docker:

```bash
docker build -t clipbridge-server .
docker run --rm -p 8787:8787 \
  -v "$(pwd)/data:/app/data" \
  -v "$(pwd)/config.yaml:/app/config.yaml:ro" \
  clipbridge-server
```

Or use the provided compose example:

```bash
cp docker-compose.example.yml docker-compose.yml
docker compose up -d
```

## Linux With systemd

Example unit:

- `deploy/systemd/clipbridge-server.service`

Typical install flow:

```bash
sudo useradd --system --home /opt/clipbridge --shell /usr/sbin/nologin clipbridge
sudo mkdir -p /opt/clipbridge/data
sudo cp clipbridge-server /opt/clipbridge/clipbridge-server
sudo cp config.yaml /opt/clipbridge/config.yaml
sudo cp deploy/systemd/clipbridge-server.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now clipbridge-server
```

## macOS With launchd

Example plist:

- `deploy/launchd/com.cinmou.clipbridge-server.plist`

Typical install flow:

```bash
sudo mkdir -p /usr/local/etc/clipbridge /usr/local/var/clipbridge/data /usr/local/var/log
sudo cp clipbridge-server /usr/local/bin/clipbridge-server
sudo cp config.yaml /usr/local/etc/clipbridge/config.yaml
cp deploy/launchd/com.cinmou.clipbridge-server.plist ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/com.cinmou.clipbridge-server.plist
```

If you use Homebrew, the example formula also includes a `service` block.

## Windows With NSSM

Example install script:

- `deploy/windows/nssm-install.ps1`

Expected layout:

```text
C:\ClipBridge\
  clipbridge-server-windows-amd64.exe
  config.yaml
  data\
```

Then install with NSSM:

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\windows\nssm-install.ps1
```

## OpenWrt

Example init script:

- `deploy/openwrt/clipbridge-server.init`

Suggested layout:

- binary at `/usr/bin/clipbridge-server`
- config at `/etc/clipbridge/config.yaml`
- working directory at `/var/lib/clipbridge`

Enable it with:

```sh
chmod +x /etc/init.d/clipbridge-server
/etc/init.d/clipbridge-server enable
/etc/init.d/clipbridge-server start
```

## Homebrew Tap Notes

An example formula is included at:

- `deploy/homebrew/clipbridge-server.rb`

Before publishing a tap, replace the placeholder SHA-256 values with real
release checksums and point the tap at your release assets.

## Scoop Notes

An example manifest is included at:

- `deploy/scoop/clipbridge-server.json`

Before publishing a Scoop bucket entry, replace the placeholder hash with the
real Windows release checksum.

## Verification Checklist

After deployment, verify:

1. the process starts with `-config config.yaml`
2. `data/clipbridge.db` is created automatically
3. `GET /api/health` returns a healthy response
4. `GET /` opens the embedded Web UI
5. one authenticated clipboard request still works

## Related Files

- `Dockerfile`
- `docker-compose.example.yml`
- `scripts/build-release.sh`
- `deploy/systemd/clipbridge-server.service`
- `deploy/launchd/com.cinmou.clipbridge-server.plist`
- `deploy/windows/nssm-install.ps1`
- `deploy/openwrt/clipbridge-server.init`
- `deploy/homebrew/clipbridge-server.rb`
- `deploy/scoop/clipbridge-server.json`
