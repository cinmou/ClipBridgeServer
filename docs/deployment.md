# Deployment

ClipBridgeServer is designed to stay simple to deploy:

- one server binary
- one `config.yaml`
- one data directory

## Local Run

Typical local layout:

```text
clipbridge/
  clipbridge-server
  config.yaml
  data/
```

Build and run:

```bash
go build -o clipbridge-server ./cmd/server
./clipbridge-server -config ./config.yaml
```

Open:

```text
http://127.0.0.1:8787/
```

## LAN Run

If you want other devices on your LAN to reach the server, set:

```yaml
server:
  host: "0.0.0.0"
  port: 8787
```

Then start normally:

```bash
./clipbridge-server -config ./config.yaml
```

Example health check from another LAN device:

```bash
curl http://YOUR-LAN-IP:8787/api/health
```

For trusted LAN use, plain HTTP is acceptable. For broader access, place the
server behind HTTPS.

## Router And NAS Notes

ClipBridgeServer can run on:

- small home servers
- NAS boxes
- SBCs
- routers with enough storage and memory

Recommendations:

- keep the data directory on durable storage
- do not store large uploads on tiny router flash
- prefer SSD, HDD, NAS volume, or external storage for long-term history

If your router or NAS has built-in HTTPS reverse proxy support, that is usually
the easiest way to add secure remote access.

## Data Directory Recommendations

The data directory holds:

- SQLite database
- uploaded files and images
- generated admin token secret
- related runtime state

Suggested practice:

- keep it outside temporary directories
- back it up if the clipboard history matters to you
- use a location with enough space for file and image uploads

## HTTPS Recommendation

For remote or public access, use an HTTPS reverse proxy such as:

- Caddy
- Nginx
- Traefik
- Tailscale
- NAS or router HTTPS gateway

Built-in TLS is optional. You can also enable it directly in config if you
already have certificate and key files.

## Optional Built-In TLS

Example:

```yaml
tls:
  enabled: true
  cert_file: "/path/to/fullchain.pem"
  key_file: "/path/to/privkey.pem"
```

Then run:

```bash
./clipbridge-server -config ./config.yaml
```

## Example Commands

Build:

```bash
go build -o clipbridge-server ./cmd/server
```

Run:

```bash
./clipbridge-server -config ./config.yaml
```

Health check:

```bash
curl http://127.0.0.1:8787/api/health
```

Open Web UI:

```text
http://127.0.0.1:8787/
```
