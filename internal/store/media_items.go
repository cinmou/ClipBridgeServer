// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// CreateClipboardItemInput is the generic write model used for text, link,
// image, and file items.
type CreateClipboardItemInput struct {
	ItemType         string
	TextContent      string
	Category         string
	SourceDeviceID   string
	SourceDeviceName string
	LocalPath        string
	Filename         string
	MIMEType         string
	SHA256           string
	URL              string
	SizeBytes        int64
	ExpiresAt        string
}

// CreateLinkItemInput creates one independent link clipboard record.
type CreateLinkItemInput struct {
	URL              string
	SourceDeviceID   string
	SourceDeviceName string
	ExpiresAt        string
}

// CreateFileItemInput creates either an image or a generic file clipboard
// record while keeping only metadata in SQLite.
type CreateFileItemInput struct {
	ItemType         string
	LocalPath        string
	Filename         string
	MIMEType         string
	SHA256           string
	SizeBytes        int64
	SourceDeviceID   string
	SourceDeviceName string
	ExpiresAt        string
}

// GetLatestClipboardItem returns the newest non-deleted clipboard item across
// all supported item types.
func (s *SQLiteStore) GetLatestClipboardItem(ctx context.Context) (*ClipboardItem, error) {
	return s.getOneClipboardItem(
		ctx,
		baseClipboardItemSelect+`
		WHERE ci.deleted_at IS NULL
		ORDER BY ci.id DESC
		LIMIT 1`,
	)
}

// ListClipboardHistory returns the full active history ordered from newest to
// oldest, optionally filtered by category.
func (s *SQLiteStore) ListClipboardHistory(ctx context.Context, categoryName string) ([]ClipboardItem, error) {
	args := make([]any, 0)
	query := baseClipboardItemSelect + `
		WHERE ci.deleted_at IS NULL`

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
		return nil, fmt.Errorf("list clipboard history: %w", err)
	}
	defer rows.Close()

	items := make([]ClipboardItem, 0)
	for rows.Next() {
		item, err := scanClipboardItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan clipboard history row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate clipboard history: %w", err)
	}

	return items, nil
}

// GetClipboardItemByID returns one non-deleted clipboard item by id.
func (s *SQLiteStore) GetClipboardItemByID(ctx context.Context, id int64) (*ClipboardItem, error) {
	return s.getOneClipboardItem(
		ctx,
		baseClipboardItemSelect+`
		WHERE ci.id = ? AND ci.deleted_at IS NULL
		LIMIT 1`,
		id,
	)
}

// DeleteClipboardItem soft-deletes one clipboard item by id regardless of its
// current item type.
func (s *SQLiteStore) DeleteClipboardItem(ctx context.Context, id int64) error {
	return s.markClipboardItemDeletedAnyType(ctx, id, time.Now().UTC().Format(time.RFC3339))
}

// SetClipboardItemFavorite toggles favorite state on any active clipboard
// record.
func (s *SQLiteStore) SetClipboardItemFavorite(ctx context.Context, id int64, favorite bool) (*ClipboardItem, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	favoriteValue := 0
	if favorite {
		favoriteValue = 1
	}

	result, err := s.db.ExecContext(
		ctx,
		`UPDATE clipboard_items
		 SET is_favorite = ?, updated_at = ?
		 WHERE id = ? AND deleted_at IS NULL`,
		favoriteValue,
		now,
		id,
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

	return s.GetClipboardItemByID(ctx, id)
}

// SetClipboardItemCategory reassigns the effective category for any active
// clipboard record.
func (s *SQLiteStore) SetClipboardItemCategory(ctx context.Context, id int64, categoryName string) (*ClipboardItem, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	normalizedName := normalizeCategoryName(categoryName)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin set item category transaction: %w", err)
	}

	if err := ensureAnyClipboardItemExistsTx(ctx, tx, id); err != nil {
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

	return s.GetClipboardItemByID(ctx, id)
}

// CreateLinkItem stores one link as its own clipboard item type.
func (s *SQLiteStore) CreateLinkItem(ctx context.Context, input CreateLinkItemInput) (*ClipboardItem, error) {
	return s.createClipboardItem(ctx, CreateClipboardItemInput{
		ItemType:         "link",
		TextContent:      strings.TrimSpace(input.URL),
		Category:         "link",
		SourceDeviceID:   input.SourceDeviceID,
		SourceDeviceName: input.SourceDeviceName,
		URL:              strings.TrimSpace(input.URL),
		SizeBytes:        int64(len([]byte(strings.TrimSpace(input.URL)))),
		ExpiresAt:        input.ExpiresAt,
	})
}

// CreateFileItem stores one image or file record with file metadata in SQLite
// and the actual bytes on disk.
func (s *SQLiteStore) CreateFileItem(ctx context.Context, input CreateFileItemInput) (*ClipboardItem, error) {
	category := "file"
	if input.ItemType == "image" {
		category = "image"
	}

	return s.createClipboardItem(ctx, CreateClipboardItemInput{
		ItemType:         input.ItemType,
		Category:         category,
		SourceDeviceID:   input.SourceDeviceID,
		SourceDeviceName: input.SourceDeviceName,
		LocalPath:        input.LocalPath,
		Filename:         input.Filename,
		MIMEType:         input.MIMEType,
		SHA256:           input.SHA256,
		SizeBytes:        input.SizeBytes,
		ExpiresAt:        input.ExpiresAt,
	})
}

func (s *SQLiteStore) createClipboardItem(ctx context.Context, input CreateClipboardItemInput) (*ClipboardItem, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	metadataJSON, err := buildMetadataJSON(clipboardMetadata{
		SourceDeviceID:   strings.TrimSpace(input.SourceDeviceID),
		SourceDeviceName: strings.TrimSpace(input.SourceDeviceName),
		LocalPath:        strings.TrimSpace(input.LocalPath),
		Filename:         strings.TrimSpace(input.Filename),
		MIMEType:         strings.TrimSpace(input.MIMEType),
		SHA256:           strings.TrimSpace(input.SHA256),
		URL:              strings.TrimSpace(input.URL),
		SizeBytes:        input.SizeBytes,
	})
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin create clipboard item transaction: %w", err)
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
		input.ItemType,
		nullIfEmpty(input.TextContent),
		metadataJSON,
		input.SizeBytes,
		nullIfEmpty(strings.TrimSpace(input.ExpiresAt)),
		now,
		now,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("insert clipboard item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("read inserted clipboard item id: %w", err)
	}

	if err := s.assignCategoryTx(ctx, tx, id, input.Category, now); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create clipboard item transaction: %w", err)
	}

	return s.GetClipboardItemByID(ctx, id)
}

func ensureAnyClipboardItemExistsTx(ctx context.Context, tx *sql.Tx, id int64) error {
	var count int
	if err := tx.QueryRowContext(
		ctx,
		`SELECT COUNT(1)
		 FROM clipboard_items
		 WHERE id = ? AND deleted_at IS NULL`,
		id,
	).Scan(&count); err != nil {
		return fmt.Errorf("check clipboard item existence: %w", err)
	}

	if count == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *SQLiteStore) markClipboardItemDeletedAnyType(ctx context.Context, id int64, deletedAt string) error {
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
