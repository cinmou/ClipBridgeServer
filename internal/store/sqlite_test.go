// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpenSQLiteCreatesDatabaseAndTables(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "clipbridge.db")

	store, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	for _, tableName := range []string{
		"schema_migrations",
		"clipboard_items",
		"categories",
		"clipboard_item_categories",
		"settings",
	} {
		if !tableExists(t, db, tableName) {
			t.Fatalf("expected table %q to exist", tableName)
		}
	}
}

func tableExists(t *testing.T, db *sql.DB, tableName string) bool {
	t.Helper()

	var count int
	err := db.QueryRow(
		"SELECT COUNT(1) FROM sqlite_master WHERE type = 'table' AND name = ?",
		tableName,
	).Scan(&count)
	if err != nil {
		t.Fatalf("query sqlite_master for %q: %v", tableName, err)
	}

	return count > 0
}
