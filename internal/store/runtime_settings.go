// SPDX-License-Identifier: GPL-3.0-only

package store

import "context"

const (
	limitsSettingsKey     = "limits.policy"
	adminTokenSettingsKey = "auth.admin_token"
)

// LoadLimitsSettings returns the persisted runtime limits or seeds defaults
// from config when the database has no override yet.
func (s *SQLiteStore) LoadLimitsSettings(ctx context.Context, defaults LimitsSettings) (LimitsSettings, error) {
	var stored LimitsSettings
	found, err := s.getJSONSetting(ctx, limitsSettingsKey, &stored)
	if err != nil {
		return LimitsSettings{}, err
	}
	if found {
		return stored, nil
	}

	if err := s.SaveLimitsSettings(ctx, defaults); err != nil {
		return LimitsSettings{}, err
	}

	return defaults, nil
}

// SaveLimitsSettings persists the runtime-editable API size limits.
func (s *SQLiteStore) SaveLimitsSettings(ctx context.Context, settings LimitsSettings) error {
	return s.setJSONSetting(ctx, limitsSettingsKey, settings)
}

// LoadAdminToken returns the persisted admin token or seeds it from config.
func (s *SQLiteStore) LoadAdminToken(ctx context.Context, defaultToken string) (string, error) {
	var stored string
	found, err := s.getJSONSetting(ctx, adminTokenSettingsKey, &stored)
	if err != nil {
		return "", err
	}
	if found {
		return stored, nil
	}

	if err := s.SaveAdminToken(ctx, defaultToken); err != nil {
		return "", err
	}

	return defaultToken, nil
}

// SaveAdminToken persists the admin token used by management APIs.
func (s *SQLiteStore) SaveAdminToken(ctx context.Context, token string) error {
	return s.setJSONSetting(ctx, adminTokenSettingsKey, token)
}
