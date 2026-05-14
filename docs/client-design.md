# Client Design

ClipBridge clients should be thin companion apps for moving clipboard items between a device and ClipBridgeServer. They are not full management consoles. Server administration, cleanup, WebDAV configuration, device management, and advanced history management should remain in the embedded Web UI.

Clients should follow the visual direction in [DESIGN.md](DESIGN.md), adapted to native desktop or mobile patterns instead of copying the Web UI layout directly.

## Product Direction

The client experience should answer three questions quickly:

- Am I connected to the right server?
- What is the latest clipboard item I know about?
- Do I want to upload my clipboard or download the latest item?

The first version should favor explicit user actions over background magic. This keeps behavior predictable and avoids surprising clipboard access.

## Desktop Client MVP

The desktop MVP should include:

- `server_url` input.
- `pairing_code` input.
- Connect button.
- Connection status.
- Upload current clipboard.
- Download latest clipboard.
- Recent history.
- Open Web UI.
- Settings.
- Optional tray or menu bar integration later.

The home screen should place the current connection status and latest known clipboard item above the primary actions. Upload and download should be large, obvious actions.

## Mobile Client MVP

The mobile MVP should include:

- `server_url` input.
- `pairing_code` input, with QR pairing reserved for a later version.
- Upload clipboard manually.
- Download latest manually.
- Recent history.
- Settings.
- No background clipboard monitoring in the first version.

Mobile clients should avoid persistent background clipboard access at first. The safer default is manual upload/download, plus future platform-native share-sheet support.

## UI Principles

- Use a simple home screen focused on clipboard transfer.
- Use large primary buttons for upload and download.
- Show the latest known clipboard item clearly.
- Avoid dashboard complexity in native clients.
- Keep advanced management in the Web UI.
- Preserve ClipBridge visual identity from `docs/DESIGN.md`.
- Adapt spacing, controls, navigation, and system affordances to each platform.
- Make connection state visible before allowing clipboard actions.

## Suggested Navigation

Desktop:

- Home: latest item, upload, download, connection status.
- History: recent items only, with basic copy/download/open actions.
- Settings: server URL, device identity, sync preferences, disconnect.
- Web UI shortcut: opens the server Web UI in the browser.

Mobile:

- Home: latest item, upload, download.
- History: recent items.
- Settings: connection, preferences, disconnect.

## What Clients Should Not Do

- Do not recreate the full Web UI management console.
- Do not call admin-only APIs.
- Do not manage cleanup, storage policy, WebDAV config, or pairing-code generation.
- Do not store the admin token.
- Do not enable background clipboard monitoring by default in the first release.

