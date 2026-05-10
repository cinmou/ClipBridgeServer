// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ImportClipboardItemInput recreates one clipboard record that arrived from an
// external sync backend while preserving the original timestamps and category.
type ImportClipboardItemInput struct {
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
	CreatedAt        string
	UpdatedAt        string
	IsFavorite       bool
	SyncKey          string
}

// ImportClipboardItem writes one synced record into SQLite if it does not
// already exist locally. It uses the same table and category mapping as normal
// API writes so history, favorites, cleanup, and downloads continue to work.
func (s *SQLiteStore) ImportClipboardItem(ctx context.Context, input ImportClipboardItemInput) (*ClipboardItem, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	createdAt := strings.TrimSpace(input.CreatedAt)
	if createdAt == "" {
		createdAt = now
	}
	updatedAt := strings.TrimSpace(input.UpdatedAt)
	if updatedAt == "" {
		updatedAt = createdAt
	}
	metadataJSON, err := buildMetadataJSON(clipboardMetadata{
		SourceDeviceID:   strings.TrimSpace(input.SourceDeviceID),
		SourceDeviceName: strings.TrimSpace(input.SourceDeviceName),
		LocalPath:        strings.TrimSpace(input.LocalPath),
		Filename:         strings.TrimSpace(input.Filename),
		MIMEType:         strings.TrimSpace(input.MIMEType),
		SHA256:           strings.TrimSpace(input.SHA256),
		URL:              strings.TrimSpace(input.URL),
		SizeBytes:        input.SizeBytes,
		SyncKey:          strings.TrimSpace(input.SyncKey),
	})
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin import clipboard item transaction: %w", err)
	}

	favoriteValue := 0
	if input.IsFavorite {
		favoriteValue = 1
	}

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO clipboard_items (
			item_type,
			text_content,
			metadata_json,
			size_bytes,
			is_favorite,
			expires_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.TrimSpace(input.ItemType),
		nullIfEmpty(strings.TrimSpace(input.TextContent)),
		metadataJSON,
		input.SizeBytes,
		favoriteValue,
		nullIfEmpty(strings.TrimSpace(input.ExpiresAt)),
		createdAt,
		updatedAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("insert imported clipboard item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("read imported clipboard item id: %w", err)
	}

	categoryName := strings.TrimSpace(input.Category)
	if categoryName == "" {
		categoryName = strings.TrimSpace(input.ItemType)
	}
	if err := s.assignCategoryTx(ctx, tx, id, categoryName, updatedAt); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit import clipboard item transaction: %w", err)
	}

	return s.GetClipboardItemByID(ctx, id)
}
