// SPDX-License-Identifier: GPL-3.0-only

package store

import "fmt"

// LimitsSettings contains the runtime-editable payload size limits enforced by
// clipboard APIs.
type LimitsSettings struct {
	MinTextBytes    int `json:"min_text_bytes"`
	MaxTextBytes    int `json:"max_text_bytes"`
	MinImageBytes   int `json:"min_image_bytes"`
	MaxImageBytes   int `json:"max_image_bytes"`
	MinFileBytes    int `json:"min_file_bytes"`
	MaxFileBytes    int `json:"max_file_bytes"`
	MinLinkBytes    int `json:"min_link_bytes"`
	MaxLinkBytes    int `json:"max_link_bytes"`
	MaxRequestBytes int `json:"max_request_bytes"`
}

// AppSettings groups the editable runtime settings with startup-only values
// that the Web UI should label as requiring restart.
type AppSettings struct {
	AdminToken            string          `json:"admin_token,omitempty"`
	Limits                LimitsSettings  `json:"limits"`
	Cleanup               CleanupSettings `json:"cleanup"`
	WebDAV                WebDAVSettings  `json:"webdav"`
	Startup               StartupSettings `json:"startup"`
	RestartRequiredFields []string        `json:"restart_required_fields"`
}

// StartupSettings are shown in the Web UI for completeness, but changing them
// still requires editing config.yaml and restarting the service.
type StartupSettings struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DataDir      string `json:"data_dir"`
	DatabasePath string `json:"database_path"`
}

// Validate keeps malformed cleanup policies from reaching the database or the
// background worker.
func (s CleanupSettings) Validate() error {
	if s.TTLHours <= 0 {
		return fmt.Errorf("ttl_hours must be greater than 0")
	}
	if s.MaxItems <= 0 {
		return fmt.Errorf("max_items must be greater than 0")
	}
	if s.MaxTotalSizeMB <= 0 {
		return fmt.Errorf("max_total_size_mb must be greater than 0")
	}
	if s.IntervalMinutes <= 0 {
		return fmt.Errorf("interval_minutes must be greater than 0")
	}
	return nil
}

// Validate keeps malformed runtime limits from breaking request validation.
func (s LimitsSettings) Validate() error {
	if s.MinTextBytes < 0 || s.MinTextBytes > s.MaxTextBytes || s.MaxTextBytes <= 0 {
		return fmt.Errorf("text byte limits are invalid")
	}
	if s.MinImageBytes < 0 || s.MinImageBytes > s.MaxImageBytes || s.MaxImageBytes <= 0 {
		return fmt.Errorf("image byte limits are invalid")
	}
	if s.MinFileBytes < 0 || s.MinFileBytes > s.MaxFileBytes || s.MaxFileBytes <= 0 {
		return fmt.Errorf("file byte limits are invalid")
	}
	if s.MinLinkBytes < 0 || s.MinLinkBytes > s.MaxLinkBytes || s.MaxLinkBytes <= 0 {
		return fmt.Errorf("link byte limits are invalid")
	}
	if s.MaxRequestBytes <= 0 {
		return fmt.Errorf("max_request_bytes must be greater than 0")
	}
	return nil
}

// WebDAVSettings describes the persisted self-hosted sync backend settings.
// They are runtime-editable because the server should not need a restart just
// to test a new endpoint or switch credentials.
type WebDAVSettings struct {
	Enabled  bool   `json:"enabled"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	BasePath string `json:"base_path"`
}

// WebDAVSyncStatus stores the latest manual sync result so the Web UI can show
// what happened without tailing server logs.
type WebDAVSyncStatus struct {
	LastSyncAt       string `json:"last_sync_at"`
	LastSuccessAt    string `json:"last_success_at"`
	LastError        string `json:"last_error,omitempty"`
	LastMessage      string `json:"last_message,omitempty"`
	TestedAt         string `json:"tested_at,omitempty"`
	LastTestSuccess  bool   `json:"last_test_success"`
	LastTestError    string `json:"last_test_error,omitempty"`
	PushedItems      int    `json:"pushed_items"`
	PulledItems      int    `json:"pulled_items"`
	PushedFiles      int    `json:"pushed_files"`
	PulledFiles      int    `json:"pulled_files"`
	RemoteItemCount  int    `json:"remote_item_count"`
	LocalItemCount   int    `json:"local_item_count"`
	ConflictSkips    int    `json:"conflict_skips"`
}

// Validate rejects obviously malformed WebDAV settings before we attempt a
// network call or persist a broken configuration.
func (s WebDAVSettings) Validate() error {
	if !s.Enabled && s.URL == "" && s.Username == "" && s.Password == "" && s.BasePath == "" {
		return nil
	}
	if s.URL == "" {
		return fmt.Errorf("url must not be empty")
	}
	if s.Username == "" {
		return fmt.Errorf("username must not be empty")
	}
	if s.Password == "" {
		return fmt.Errorf("password must not be empty")
	}
	return nil
}
