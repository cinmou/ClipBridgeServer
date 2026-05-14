# Security

ClipBridgeServer `Cherwell` uses a light, practical security model for
self-hosted clipboard use. It is meant for personal servers, home labs, and
trusted LAN deployments rather than enterprise zero-trust environments.

## Current Security Model

- Default bind address is `127.0.0.1`
- If `auth.token` is empty, the server generates an admin token automatically
- The generated admin token is stored at `data/secrets/admin_token`
- Pairing codes are short-lived
- Pairing codes are single-use
- Device tokens are generated during pairing
- Trusted LAN HTTP is acceptable for local use
- Remote or public access should use an HTTPS reverse proxy
- Built-in TLS is optional for advanced users
- Logs should not expose tokens, pairing codes, cookies, or WebDAV passwords

## Admin Token

The admin token is the main management credential for:

- device pairing code generation
- device revocation
- settings access
- cleanup and storage status
- WebDAV settings and sync preview

If you leave `auth.token` empty in `config.yaml`, ClipBridgeServer will:

1. generate a token on first startup
2. save it locally under `data/secrets/admin_token`
3. reuse that token on later restarts

This avoids weak baked-in defaults while keeping deployment simple.

## Pairing And Device Access

ClipBridgeServer separates management access from clipboard client access:

- admin token: management APIs
- device token: clipboard APIs after pairing

Current pairing behavior:

- pairing codes expire after 5 minutes
- pairing codes become invalid after one successful use
- device tokens are created during pairing and then reused by the device

## LAN And Remote Access

For a personal server on a trusted network, plain HTTP is acceptable if you
understand the network boundary and do not expose the service to the public
internet.

Recommended guidance:

- local only: keep `server.host = 127.0.0.1`
- trusted LAN: you may use `server.host = 0.0.0.0` over HTTP
- remote or public access: use HTTPS through a reverse proxy

Typical HTTPS reverse proxy choices:

- Caddy
- Nginx
- Traefik
- Tailscale
- NAS or router HTTPS gateway

Built-in TLS is available, but it is optional and intended for users who
already manage certificate files directly.

## Logging

ClipBridgeServer should not log raw values for:

- admin tokens
- device tokens
- pairing codes
- cookies
- WebDAV passwords

This reduces accidental credential leakage in normal operation and debugging.

## Future Hardening

The current model is intentionally small and practical. Stronger security for
external cloud or WebDAV sync workflows may be added later, but it is not
claimed as complete in the current beta.
