// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// ErrNotFound lets higher layers decide whether a missing record maps to HTTP
// 404 or another API-level response.
var ErrNotFound = errors.New("store: record not found")

// ErrConflict marks writes that would violate a uniqueness rule, such as
// creating a category with a name that already exists.
var ErrConflict = errors.New("store: conflict")

// ClipboardItem is the storage representation returned to the API layer.
// Keeping this structure explicit makes the Web UI and API handlers easier to
// evolve as more clipboard item types arrive in later phases.
type ClipboardItem struct {
	ID               int64
	ItemType         string
	TextContent      string
	MetadataJSON     string
	IsFavorite       bool
	Category         string
	SourceDeviceID   string
	SourceDeviceName string
	LocalPath        string
	SizeBytes        int64
	ExpiresAt        string
	CreatedAt        string
	UpdatedAt        string
}

// CreateTextItemInput keeps text creation extensible without introducing a
// dedicated Web UI-only API surface.
type CreateTextItemInput struct {
	Text             string
	SourceDeviceID   string
	SourceDeviceName string
	ExpiresAt        string
}

// Category is one user-visible category that can be attached to clipboard
// items. The first wave includes built-ins like text/image/link/file, but we
// keep the storage generic so custom categories fit the same model.
type Category struct {
	ID        int64
	Name      string
	CreatedAt string
	UpdatedAt string
}

// CleanupSettings defines the persisted retention policy that the background
// cleaner and admin API both use.
type CleanupSettings struct {
	TTLHours        int  `json:"ttl_hours"`
	MaxItems        int  `json:"max_items"`
	MaxTotalSizeMB  int  `json:"max_total_size_mb"`
	IntervalMinutes int  `json:"interval_minutes"`
	Enabled         bool `json:"enabled"`
}

// CleanupStatus stores the most recent worker execution summary so the Web UI
// can explain what happened without reading logs.
type CleanupStatus struct {
	LastRunAt           string `json:"last_run_at"`
	LastRunReason       string `json:"last_run_reason"`
	DeletedExpired      int    `json:"deleted_expired"`
	DeletedMaxItems     int    `json:"deleted_max_items"`
	DeletedStorage      int    `json:"deleted_storage"`
	DeletedFiles        int    `json:"deleted_files"`
	LastError           string `json:"last_error,omitempty"`
	HistoryCount        int    `json:"history_count"`
	FavoriteCount       int    `json:"favorite_count"`
	TotalBytes          int64  `json:"total_bytes"`
	FileBytes           int64  `json:"file_bytes"`
	NonFavoriteFileSize int64  `json:"non_favorite_file_size"`
}

// StorageStatus is a lightweight point-in-time summary of retained data.
type StorageStatus struct {
	HistoryCount  int   `json:"history_count"`
	FavoriteCount int   `json:"favorite_count"`
	TotalBytes    int64 `json:"total_bytes"`
	FileBytes     int64 `json:"file_bytes"`
}

// CleanupCandidate is the minimal record shape the cleaner needs when deciding
// what can be deleted under retention pressure.
type CleanupCandidate struct {
	ID         int64
	ItemType   string
	IsFavorite bool
	Metadata   clipboardMetadata
	SizeBytes  int64
	ExpiresAt  string
	CreatedAt  string
}

type clipboardMetadata struct {
	SourceDeviceID   string `json:"source_device_id,omitempty"`
	SourceDeviceName string `json:"source_device_name,omitempty"`
	LocalPath        string `json:"local_path,omitempty"`
	SizeBytes        int64  `json:"size_bytes,omitempty"`
}

const (
	cleanupSettingsKey = "cleanup.policy"
	cleanupStatusKey   = "cleanup.status"
)

// SQLiteStore owns the one SQLite connection pool used by the server process.
// The service is intentionally a single-binary deployment, so the store keeps
// all persistence behavior inside the same executable.
type SQLiteStore struct {
	db *sql.DB
}

// OpenSQLite creates the local data directory, opens the database, enables
// foreign keys, and runs schema migrations before the HTTP server starts.
func OpenSQLite(databasePath string) (*SQLiteStore, error) {
	if err := ensureDatabaseDir(databasePath); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database %q: %w", databasePath, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable sqlite foreign keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database %q: %w", databasePath, err)
	}

	if err := RunMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

// Close releases the SQLite connection pool on process shutdown.
func (s *SQLiteStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	return s.db.Close()
}

// CreateTextItem writes one text clipboard record, then attaches the built-in
// "text" category so history filtering works consistently with older and newer
// items alike.
func (s *SQLiteStore) CreateTextItem(ctx context.Context, input CreateTextItemInput) (*ClipboardItem, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	textBytes := int64(len([]byte(input.Text)))
	metadataJSON, err := buildMetadataJSON(clipboardMetadata{
		SourceDeviceID:   strings.TrimSpace(input.SourceDeviceID),
		SourceDeviceName: strings.TrimSpace(input.SourceDeviceName),
		SizeBytes:        textBytes,
	})
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin create text item transaction: %w", err)
	}

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO clipboard_items (
			item_type,
			text_content,
			metadata_json,
			size_bytes,
			expires_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"text",
		input.Text,
		metadataJSON,
		textBytes,
		nullIfEmpty(strings.TrimSpace(input.ExpiresAt)),
		now,
		now,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("insert text clipboard item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("read inserted clipboard item id: %w", err)
	}

	if err := s.assignCategoryTx(ctx, tx, id, "text", now); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create text item transaction: %w", err)
	}

	return s.GetTextItemByID(ctx, id)
}

// GetLatestTextItem returns the newest non-deleted text item.
func (s *SQLiteStore) GetLatestTextItem(ctx context.Context) (*ClipboardItem, error) {
	return s.getOneClipboardItem(
		ctx,
		baseClipboardItemSelect+`
		WHERE ci.item_type = ? AND ci.deleted_at IS NULL
		ORDER BY ci.id DESC
		LIMIT 1`,
		"text",
	)
}

// ListTextHistory returns non-deleted text clipboard items ordered from newest
// to oldest. When categoryName is provided, the result is filtered to that one
// category, which lets the Web UI and API build simple history filters without
// changing the response shape.
func (s *SQLiteStore) ListTextHistory(ctx context.Context, categoryName string) ([]ClipboardItem, error) {
	args := []any{"text"}
	query := baseClipboardItemSelect + `
		WHERE ci.item_type = ? AND ci.deleted_at IS NULL`

	if categoryName != "" {
		query += ` AND EXISTS (
			SELECT 1
			FROM clipboard_item_categories AS filter_cic
			JOIN categories AS filter_c ON filter_c.id = filter_cic.category_id
			WHERE filter_cic.clipboard_item_id = ci.id
			  AND filter_c.name = ?
		)`
		args = append(args, strings.TrimSpace(categoryName))
	}

	query += `
		ORDER BY ci.id DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list text clipboard history: %w", err)
	}
	defer rows.Close()

	items := make([]ClipboardItem, 0)
	for rows.Next() {
		item, err := scanClipboardItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan text clipboard history row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate text clipboard history: %w", err)
	}

	return items, nil
}

// GetTextItemByID returns one non-deleted text clipboard record by id.
func (s *SQLiteStore) GetTextItemByID(ctx context.Context, id int64) (*ClipboardItem, error) {
	return s.getOneClipboardItem(
		ctx,
		baseClipboardItemSelect+`
		WHERE ci.id = ? AND ci.item_type = ? AND ci.deleted_at IS NULL
		LIMIT 1`,
		id,
		"text",
	)
}

// DeleteTextItem performs a soft delete so later cleanup phases can respect
// favorites and retention policies without losing deletion timestamps.
func (s *SQLiteStore) DeleteTextItem(ctx context.Context, id int64) error {
	return s.markClipboardItemDeleted(ctx, id, "text", time.Now().UTC().Format(time.RFC3339))
}

// SetFavorite toggles whether one text clipboard record is pinned as a
// favorite. Favorites are the records a future TTL worker must always preserve.
func (s *SQLiteStore) SetFavorite(ctx context.Context, id int64, favorite bool) (*ClipboardItem, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	favoriteValue := 0
	if favorite {
		favoriteValue = 1
	}

	result, err := s.db.ExecContext(
		ctx,
		`UPDATE clipboard_items
		 SET is_favorite = ?, updated_at = ?
		 WHERE id = ? AND item_type = ? AND deleted_at IS NULL`,
		favoriteValue,
		now,
		id,
		"text",
	)
	if err != nil {
		return nil, fmt.Errorf("update favorite state for item %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("read updated favorite row count: %w", err)
	}

	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return s.GetTextItemByID(ctx, id)
}

// ListFavorites returns all non-deleted favorite clipboard items ordered from
// newest to oldest.
func (s *SQLiteStore) ListFavorites(ctx context.Context) ([]ClipboardItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		baseClipboardItemSelect+`
		WHERE ci.item_type = ? AND ci.deleted_at IS NULL AND ci.is_favorite = 1
		ORDER BY ci.id DESC`,
		"text",
	)
	if err != nil {
		return nil, fmt.Errorf("list favorites: %w", err)
	}
	defer rows.Close()

	items := make([]ClipboardItem, 0)
	for rows.Next() {
		item, err := scanClipboardItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan favorite row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate favorites: %w", err)
	}

	return items, nil
}

// ListCategories returns every known category, including the built-ins seeded
// by migrations and any custom categories users add later.
func (s *SQLiteStore) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, name, created_at, updated_at
		 FROM categories
		 ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan category row: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}

	return categories, nil
}

// CreateCategory adds one new custom category.
func (s *SQLiteStore) CreateCategory(ctx context.Context, name string) (*Category, error) {
	normalizedName := normalizeCategoryName(name)
	now := time.Now().UTC().Format(time.RFC3339)

	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO categories (name, created_at, updated_at)
		 VALUES (?, ?, ?)`,
		normalizedName,
		now,
		now,
	)
	if err != nil {
		if isUniqueConstraint(err) {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("create category %q: %w", normalizedName, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read inserted category id: %w", err)
	}

	return s.getOneCategory(
		ctx,
		`SELECT id, name, created_at, updated_at
		 FROM categories
		 WHERE id = ?
		 LIMIT 1`,
		id,
	)
}

// SetItemCategory replaces the current category assignment for one text record.
// We keep one effective category per item in this phase so the Web UI can offer
// a simple filter and edit model.
func (s *SQLiteStore) SetItemCategory(ctx context.Context, id int64, categoryName string) (*ClipboardItem, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	normalizedName := normalizeCategoryName(categoryName)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin set item category transaction: %w", err)
	}

	if err := ensureClipboardItemExistsTx(ctx, tx, id); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := s.assignCategoryTx(ctx, tx, id, normalizedName, now); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE clipboard_items
		 SET updated_at = ?
		 WHERE id = ?`,
		now,
		id,
	); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("update item category timestamp: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit set item category transaction: %w", err)
	}

	return s.GetTextItemByID(ctx, id)
}

// LoadCleanupSettings returns the persisted cleanup policy or seeds the
// defaults from config when the database has not stored an override yet.
func (s *SQLiteStore) LoadCleanupSettings(ctx context.Context, defaults CleanupSettings) (CleanupSettings, error) {
	var stored CleanupSettings
	found, err := s.getJSONSetting(ctx, cleanupSettingsKey, &stored)
	if err != nil {
		return CleanupSettings{}, err
	}
	if found {
		return stored, nil
	}

	if err := s.SaveCleanupSettings(ctx, defaults); err != nil {
		return CleanupSettings{}, err
	}

	return defaults, nil
}

// SaveCleanupSettings persists the currently active retention policy.
func (s *SQLiteStore) SaveCleanupSettings(ctx context.Context, settings CleanupSettings) error {
	return s.setJSONSetting(ctx, cleanupSettingsKey, settings)
}

// LoadCleanupStatus returns the latest cleaner run summary.
func (s *SQLiteStore) LoadCleanupStatus(ctx context.Context) (CleanupStatus, error) {
	var status CleanupStatus
	found, err := s.getJSONSetting(ctx, cleanupStatusKey, &status)
	if err != nil {
		return CleanupStatus{}, err
	}
	if !found {
		return CleanupStatus{}, nil
	}

	return status, nil
}

// SaveCleanupStatus persists the latest cleaner run summary.
func (s *SQLiteStore) SaveCleanupStatus(ctx context.Context, status CleanupStatus) error {
	return s.setJSONSetting(ctx, cleanupStatusKey, status)
}

// RefreshExpirations keeps non-favorite clipboard records aligned with the
// current TTL policy so cleanup logic can work from one explicit expires_at
// field instead of recalculating at every deletion point.
func (s *SQLiteStore) RefreshExpirations(ctx context.Context, ttlHours int) error {
	modifier := fmt.Sprintf("+%d hours", ttlHours)
	if _, err := s.db.ExecContext(
		ctx,
		`UPDATE clipboard_items
		 SET expires_at = CASE
		 	WHEN is_favorite = 1 OR deleted_at IS NOT NULL THEN NULL
		 	ELSE STRFTIME('%Y-%m-%dT%H:%M:%SZ', DATETIME(created_at, ?))
		 END`,
		modifier,
	); err != nil {
		return fmt.Errorf("refresh clipboard expirations: %w", err)
	}

	return nil
}

// ListCleanupCandidates returns every active clipboard item sorted from oldest
// to newest so the cleaner can remove stale data in deterministic order.
func (s *SQLiteStore) ListCleanupCandidates(ctx context.Context) ([]CleanupCandidate, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, item_type, is_favorite, COALESCE(metadata_json, '{}'),
		        COALESCE(size_bytes, LENGTH(COALESCE(text_content, ''))),
		        COALESCE(expires_at, ''), created_at
		 FROM clipboard_items
		 WHERE deleted_at IS NULL
		 ORDER BY created_at ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list cleanup candidates: %w", err)
	}
	defer rows.Close()

	candidates := make([]CleanupCandidate, 0)
	for rows.Next() {
		var candidate CleanupCandidate
		var favoriteValue int
		var metadataJSON string

		if err := rows.Scan(
			&candidate.ID,
			&candidate.ItemType,
			&favoriteValue,
			&metadataJSON,
			&candidate.SizeBytes,
			&candidate.ExpiresAt,
			&candidate.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan cleanup candidate: %w", err)
		}

		candidate.IsFavorite = favoriteValue == 1
		candidate.Metadata = parseClipboardMetadata(metadataJSON, candidate.SizeBytes)
		candidates = append(candidates, candidate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cleanup candidates: %w", err)
	}

	return candidates, nil
}

// MarkItemDeletedForCleanup performs the same soft-delete strategy as manual
// deletion, but it works across future clipboard item types.
func (s *SQLiteStore) MarkItemDeletedForCleanup(ctx context.Context, id int64, deletedAt string) error {
	result, err := s.db.ExecContext(
		ctx,
		`UPDATE clipboard_items
		 SET deleted_at = ?, updated_at = ?
		 WHERE id = ? AND deleted_at IS NULL`,
		deletedAt,
		deletedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft delete clipboard item %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read cleanup delete row count: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetStorageStatus summarizes active clipboard storage pressure for the admin
// API and embedded Web UI.
func (s *SQLiteStore) GetStorageStatus(ctx context.Context) (StorageStatus, error) {
	var status StorageStatus
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT
			COUNT(1),
			COALESCE(SUM(CASE WHEN is_favorite = 1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(size_bytes), 0),
			COALESCE(SUM(CASE WHEN item_type = 'file' THEN size_bytes ELSE 0 END), 0)
		 FROM clipboard_items
		 WHERE deleted_at IS NULL`,
	).Scan(
		&status.HistoryCount,
		&status.FavoriteCount,
		&status.TotalBytes,
		&status.FileBytes,
	); err != nil {
		return StorageStatus{}, fmt.Errorf("query storage status: %w", err)
	}

	return status, nil
}

const baseClipboardItemSelect = `
SELECT
	ci.id,
	ci.item_type,
	COALESCE(ci.text_content, ''),
	COALESCE(ci.metadata_json, '{}'),
	ci.is_favorite,
	COALESCE((
		SELECT c.name
		FROM clipboard_item_categories AS cic
		JOIN categories AS c ON c.id = cic.category_id
		WHERE cic.clipboard_item_id = ci.id
		ORDER BY cic.created_at DESC, c.id DESC
		LIMIT 1
	), ''),
	COALESCE(ci.size_bytes, LENGTH(COALESCE(ci.text_content, ''))),
	COALESCE(ci.expires_at, ''),
	ci.created_at,
	ci.updated_at
FROM clipboard_items AS ci`

func (s *SQLiteStore) getOneClipboardItem(ctx context.Context, query string, args ...any) (*ClipboardItem, error) {
	row := s.db.QueryRowContext(ctx, query, args...)
	item, err := scanClipboardItem(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query clipboard item: %w", err)
	}

	return &item, nil
}

func scanClipboardItem(scanner interface {
	Scan(dest ...any) error
}) (ClipboardItem, error) {
	var item ClipboardItem
	var favoriteValue int
	var metadataJSON string

	err := scanner.Scan(
		&item.ID,
		&item.ItemType,
		&item.TextContent,
		&metadataJSON,
		&favoriteValue,
		&item.Category,
		&item.SizeBytes,
		&item.ExpiresAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return ClipboardItem{}, err
	}

	item.MetadataJSON = metadataJSON
	item.IsFavorite = favoriteValue == 1

	metadata := parseClipboardMetadata(metadataJSON, item.SizeBytes)
	item.SourceDeviceID = metadata.SourceDeviceID
	item.SourceDeviceName = metadata.SourceDeviceName
	item.LocalPath = metadata.LocalPath

	return item, nil
}

func (s *SQLiteStore) assignCategoryTx(ctx context.Context, tx *sql.Tx, itemID int64, categoryName string, createdAt string) error {
	category, err := s.getOneCategoryTx(
		ctx,
		tx,
		`SELECT id, name, created_at, updated_at
		 FROM categories
		 WHERE name = ?
		 LIMIT 1`,
		categoryName,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fmt.Errorf("category %q not found: %w", categoryName, ErrNotFound)
		}
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`DELETE FROM clipboard_item_categories
		 WHERE clipboard_item_id = ?`,
		itemID,
	); err != nil {
		return fmt.Errorf("clear existing item categories: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO clipboard_item_categories (clipboard_item_id, category_id, created_at)
		 VALUES (?, ?, ?)`,
		itemID,
		category.ID,
		createdAt,
	); err != nil {
		return fmt.Errorf("assign category %q to item %d: %w", categoryName, itemID, err)
	}

	return nil
}

func ensureClipboardItemExistsTx(ctx context.Context, tx *sql.Tx, id int64) error {
	var count int
	if err := tx.QueryRowContext(
		ctx,
		`SELECT COUNT(1)
		 FROM clipboard_items
		 WHERE id = ? AND item_type = ? AND deleted_at IS NULL`,
		id,
		"text",
	).Scan(&count); err != nil {
		return fmt.Errorf("check clipboard item existence: %w", err)
	}

	if count == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *SQLiteStore) getOneCategory(ctx context.Context, query string, args ...any) (*Category, error) {
	var category Category

	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query category: %w", err)
	}

	return &category, nil
}

func (s *SQLiteStore) getOneCategoryTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (*Category, error) {
	var category Category

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query category in transaction: %w", err)
	}

	return &category, nil
}

func (s *SQLiteStore) getJSONSetting(ctx context.Context, key string, target any) (bool, error) {
	var rawValue string
	err := s.db.QueryRowContext(
		ctx,
		`SELECT value
		 FROM settings
		 WHERE key = ?
		 LIMIT 1`,
		key,
	).Scan(&rawValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("query setting %q: %w", key, err)
	}

	if err := json.Unmarshal([]byte(rawValue), target); err != nil {
		return false, fmt.Errorf("decode setting %q: %w", key, err)
	}

	return true, nil
}

func (s *SQLiteStore) setJSONSetting(ctx context.Context, key string, value any) error {
	rawValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode setting %q: %w", key, err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := s.db.ExecContext(
		ctx,
		`INSERT INTO settings (key, value, updated_at)
		 VALUES (?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		 	value = excluded.value,
		 	updated_at = excluded.updated_at`,
		key,
		string(rawValue),
		now,
	); err != nil {
		return fmt.Errorf("save setting %q: %w", key, err)
	}

	return nil
}

func (s *SQLiteStore) markClipboardItemDeleted(ctx context.Context, id int64, itemType string, deletedAt string) error {
	result, err := s.db.ExecContext(
		ctx,
		`UPDATE clipboard_items
		 SET deleted_at = ?, updated_at = ?
		 WHERE id = ? AND item_type = ? AND deleted_at IS NULL`,
		deletedAt,
		deletedAt,
		id,
		itemType,
	)
	if err != nil {
		return fmt.Errorf("delete clipboard item %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted clipboard row count: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func buildMetadataJSON(metadata clipboardMetadata) (string, error) {
	rawValue, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("encode clipboard metadata: %w", err)
	}

	return string(rawValue), nil
}

func parseClipboardMetadata(rawValue string, fallbackSizeBytes int64) clipboardMetadata {
	metadata := clipboardMetadata{SizeBytes: fallbackSizeBytes}
	if strings.TrimSpace(rawValue) == "" {
		return metadata
	}

	if err := json.Unmarshal([]byte(rawValue), &metadata); err != nil {
		return clipboardMetadata{SizeBytes: fallbackSizeBytes}
	}

	if metadata.SizeBytes <= 0 {
		metadata.SizeBytes = fallbackSizeBytes
	}

	return metadata
}

func normalizeCategoryName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return value
}

func ensureDatabaseDir(databasePath string) error {
	dir := filepath.Dir(databasePath)
	if dir == "." || dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create sqlite data directory %q: %w", dir, err)
	}

	return nil
}
