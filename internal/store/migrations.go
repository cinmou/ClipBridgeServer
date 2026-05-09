// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var migrationVersions = []string{
	"001_init.sql",
	"002_device_pairing.sql",
	"003_favorites_categories.sql",
	"004_cleanup_metadata.sql",
}

func RunMigrations(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	for _, version := range migrationVersions {
		applied, err := hasMigration(db, version)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		if err := applyMigration(db, version); err != nil {
			return err
		}
	}

	return nil
}

func applyMigration(db *sql.DB, version string) error {
	script, err := loadMigration(version)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}

	if _, err := tx.Exec(script); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("apply migration %s: %w", version, err)
	}

	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version) VALUES (?)",
		version,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %s: %w", version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", version, err)
	}

	return nil
}

func ensureMigrationsTable(db *sql.DB) error {
	const sqlText = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version TEXT PRIMARY KEY,
	applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	if _, err := db.Exec(sqlText); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	return nil
}

func hasMigration(db *sql.DB, version string) (bool, error) {
	var count int

	if err := db.QueryRow(
		"SELECT COUNT(1) FROM schema_migrations WHERE version = ?",
		version,
	).Scan(&count); err != nil {
		return false, fmt.Errorf("query schema_migrations for %s: %w", version, err)
	}

	return count > 0, nil
}

func loadMigration(version string) (string, error) {
	path, err := migrationPath(version)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read migration %q: %w", path, err)
	}

	return string(data), nil
}

func migrationPath(version string) (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve migration path: runtime caller unavailable")
	}

	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	return filepath.Join(repoRoot, "migrations", version), nil
}
