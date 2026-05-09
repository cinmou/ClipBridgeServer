// SPDX-License-Identifier: GPL-3.0-only

package config

func Default() *Config {
	return &Config{
		Auth: AuthConfig{
			Token: "dev-token-please-change",
		},
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 8787,
		},
		Storage: StorageConfig{
			DataDir:        "./data",
			DatabasePath:   "./data/clipbridge.db",
			TTLHours:       168,
			MaxItems:       1000,
			MaxTotalSizeMB: 2048,
		},
		Limits: LimitsConfig{
			MinTextBytes:    1,
			MaxTextBytes:    262144,
			MaxRequestBytes: 1048576,
		},
		Cleaner: CleanerConfig{
			Enabled:         true,
			IntervalMinutes: 30,
		},
	}
}
