# Web UI

ClipBridgeServer ships with an embedded Web UI served directly from the server
binary at `GET /`.

The current Cherwell beta Web UI supports English and Simplified Chinese.

## Dashboard / Overview

The Dashboard is for clipboard actions and quick status:

- latest clipboard item
- quick text upload
- fetch latest item
- copy latest text or link into the browser clipboard
- high-level service summary

Access Token controls are not on the Dashboard.

## History

The History page is for browsing and managing clipboard items:

- reverse chronological list
- search
- type filters
- item details
- copy, open, download, favorite, unfavorite, delete

## Favorites

The Favorites page shows all favorite items.

Favorites are protected from automatic cleanup until you remove the favorite
flag.

## Devices / Pairing

The Devices page is for device onboarding and trust management:

- generate pairing codes
- view expiry information
- view paired devices
- revoke devices

Pairing codes are short-lived and single-use.

## Settings

The Settings page contains:

- Access Token controls
- language switcher
- security and access information
- runtime details
- cleanup policy
- storage usage
- cleanup status

Settings shows security/access information such as bind mode, TLS status, and
pairing policy.

## WebDAV

The WebDAV page is a preview/manual sync page:

- save WebDAV settings
- test WebDAV connection
- run one manual sync
- view latest sync state

WebDAV sync is currently preview-oriented and manual rather than always-on.

## Language Switcher

The Web UI supports:

- English
- Simplified Chinese (`zh-CN`)

Language selection is available in Settings and is stored in the browser.
