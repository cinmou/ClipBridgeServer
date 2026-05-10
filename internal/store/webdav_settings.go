// SPDX-License-Identifier: GPL-3.0-only

package store

import "context"

const (
	webDAVSettingsKey = "webdav.settings"
	webDAVStatusKey   = "webdav.status"
)

// LoadWebDAVSettings returns the persisted WebDAV configuration or a blank
// default when the feature has not been configured yet.
func (s *SQLiteStore) LoadWebDAVSettings(ctx context.Context) (WebDAVSettings, error) {
	var stored WebDAVSettings
	found, err := s.getJSONSetting(ctx, webDAVSettingsKey, &stored)
	if err != nil {
		return WebDAVSettings{}, err
	}
	if !found {
		return WebDAVSettings{BasePath: "ClipBridgeServer"}, nil
	}
	if stored.BasePath == "" {
		stored.BasePath = "ClipBridgeServer"
	}
	return stored, nil
}

// SaveWebDAVSettings persists the editable WebDAV backend configuration.
func (s *SQLiteStore) SaveWebDAVSettings(ctx context.Context, settings WebDAVSettings) error {
	return s.setJSONSetting(ctx, webDAVSettingsKey, settings)
}

// LoadWebDAVSyncStatus returns the latest saved sync status, or an empty value
// before the first connection test or manual sync.
func (s *SQLiteStore) LoadWebDAVSyncStatus(ctx context.Context) (WebDAVSyncStatus, error) {
	var status WebDAVSyncStatus
	found, err := s.getJSONSetting(ctx, webDAVStatusKey, &status)
	if err != nil {
		return WebDAVSyncStatus{}, err
	}
	if !found {
		return WebDAVSyncStatus{}, nil
	}
	return status, nil
}

// SaveWebDAVSyncStatus persists the latest connection test or manual sync
// result for the Web UI and admin API.
func (s *SQLiteStore) SaveWebDAVSyncStatus(ctx context.Context, status WebDAVSyncStatus) error {
	return s.setJSONSetting(ctx, webDAVStatusKey, status)
}
