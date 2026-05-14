# Client Storage

ClipBridge clients should store only the local state needed to reconnect, sync manually, and present a small recent clipboard view. The admin token should never be stored by a client.

## Required Local State

Clients should store:

- `server_url`
- `device_id`
- `device_token`
- `device_name`
- `last_sync_at`
- `last_seen_item_id`
- User preferences
- `auto_sync_enabled`
- `upload_text_enabled`
- `upload_file_enabled`

## Sensitive Values

The `device_token` is sensitive. Store it in platform secure storage where available:

- macOS/iOS: Keychain.
- Windows: Credential Manager or DPAPI-backed storage.
- Android: Keystore-backed encrypted storage.
- Linux: Secret Service, KWallet, or another desktop secret store when available.

If secure storage is unavailable, clients should make that limitation clear in documentation and avoid storing broader privileges than the device token.

## Non-Sensitive Preferences

The following values may be stored in normal app preferences:

- `server_url`
- `device_name`
- `last_sync_at`
- `last_seen_item_id`
- Language or theme preference.
- Upload/download behavior preferences.

Client implementations may still choose to store `server_url` with credentials for simplicity, but the token itself should remain protected.

## Cached Clipboard Metadata

Clients may cache recent item metadata for responsiveness:

- Item ID.
- Type.
- Preview text.
- Created time.
- MIME type.
- Size.
- Source device name.

Clients should avoid caching file payloads by default unless the user explicitly downloads them or the platform requires a temporary file.

## Disconnect Behavior

Every client should provide a clear disconnect action.

Disconnect should:

- Delete `device_token` from local storage.
- Delete `device_id` if it is only useful with the token.
- Clear cached recent history if the user chooses.
- Keep harmless preferences such as theme or language unless the user requests a full reset.

Disconnecting locally does not necessarily revoke the device on the server. Device revocation should remain an admin action in the Web UI.

## What Not To Store

Clients should not store:

- Admin token.
- Pairing codes after pairing completes.
- WebDAV credentials.
- Server cleanup or storage configuration.
- Full clipboard history unless the user explicitly enables local caching.

