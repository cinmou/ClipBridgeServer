# Client Roadmap

ClipBridge clients are future companion apps. They do not exist as completed products in the current Cherwell beta.

## 1. Client MVP

Goal: a thin manual client that can pair with ClipBridgeServer and move clipboard items on demand.

Candidate scope:

- Connect with `server_url` and `pairing_code`.
- Store a `device_token`.
- Upload current clipboard manually.
- Download latest clipboard manually.
- Show latest known item.
- Show a short recent-history list.
- Open the embedded Web UI for management.
- Disconnect and clear local credentials.

## 2. Desktop Tray Or Menu Bar

Goal: make common desktop actions quick without turning the client into a management console.

Candidate scope:

- Tray/menu bar status.
- Upload current clipboard.
- Download latest clipboard.
- Open recent items.
- Open Web UI.
- Optional background sync after the manual MVP is stable.

## 3. Mobile Manual Client

Goal: provide a safe mobile companion without background clipboard monitoring in the first version.

Candidate scope:

- Pair with server URL and pairing code.
- Upload clipboard manually.
- Download latest manually.
- View recent history.
- Basic settings and disconnect.

## 4. Share Extension / Share Sheet

Goal: integrate with platform sharing flows after the manual client is useful.

Candidate scope:

- Send text or links to ClipBridge from other apps.
- Send images/files when supported.
- Open downloaded items in native share flows.

## 5. Realtime Sync Later

Goal: explore faster sync only after the API and manual clients are stable.

Possible directions:

- Polling improvements.
- Server-sent events.
- WebSocket-based updates.
- Platform-specific background sync where appropriate.

Realtime behavior should remain optional and should not compromise the simple self-hosted security model.

