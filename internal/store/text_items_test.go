// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestTextClipboardCRUD(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "clipbridge.db")
	dbStore, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	t.Cleanup(func() {
		_ = dbStore.Close()
	})

	ctx := context.Background()

	firstItem, err := dbStore.CreateTextItem(ctx, CreateTextItemInput{
		Text:             "first",
		SourceDeviceID:   "web-ui",
		SourceDeviceName: "Web UI",
	})
	if err != nil {
		t.Fatalf("CreateTextItem(first) error = %v", err)
	}

	secondItem, err := dbStore.CreateTextItem(ctx, CreateTextItemInput{
		Text:             "second",
		SourceDeviceID:   "macbook",
		SourceDeviceName: "MacBook Pro",
	})
	if err != nil {
		t.Fatalf("CreateTextItem(second) error = %v", err)
	}

	latestItem, err := dbStore.GetLatestTextItem(ctx)
	if err != nil {
		t.Fatalf("GetLatestTextItem() error = %v", err)
	}
	if latestItem.ID != secondItem.ID || latestItem.TextContent != "second" {
		t.Fatalf("GetLatestTextItem() = %+v, want second item", latestItem)
	}
	if latestItem.Category != "text" {
		t.Fatalf("GetLatestTextItem() category = %q, want %q", latestItem.Category, "text")
	}

	history, err := dbStore.ListTextHistory(ctx, "")
	if err != nil {
		t.Fatalf("ListTextHistory() error = %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("ListTextHistory() len = %d, want 2", len(history))
	}
	if history[0].ID != secondItem.ID || history[1].ID != firstItem.ID {
		t.Fatalf("ListTextHistory() order = %+v, want newest first", history)
	}

	gotItem, err := dbStore.GetTextItemByID(ctx, firstItem.ID)
	if err != nil {
		t.Fatalf("GetTextItemByID() error = %v", err)
	}
	if gotItem.TextContent != "first" {
		t.Fatalf("GetTextItemByID() text = %q, want %q", gotItem.TextContent, "first")
	}
	if gotItem.Category != "text" {
		t.Fatalf("GetTextItemByID() category = %q, want %q", gotItem.Category, "text")
	}
	if gotItem.SourceDeviceName != "Web UI" {
		t.Fatalf("GetTextItemByID() sourceDeviceName = %q, want %q", gotItem.SourceDeviceName, "Web UI")
	}

	favoritedItem, err := dbStore.SetFavorite(ctx, firstItem.ID, true)
	if err != nil {
		t.Fatalf("SetFavorite(true) error = %v", err)
	}
	if !favoritedItem.IsFavorite {
		t.Fatalf("SetFavorite(true) isFavorite = %v, want true", favoritedItem.IsFavorite)
	}

	favorites, err := dbStore.ListFavorites(ctx)
	if err != nil {
		t.Fatalf("ListFavorites() error = %v", err)
	}
	if len(favorites) != 1 || favorites[0].ID != firstItem.ID {
		t.Fatalf("ListFavorites() = %+v, want only first item", favorites)
	}

	category, err := dbStore.CreateCategory(ctx, "work")
	if err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}
	if category.Name != "work" {
		t.Fatalf("CreateCategory() name = %q, want %q", category.Name, "work")
	}

	retypedItem, err := dbStore.SetItemCategory(ctx, secondItem.ID, "work")
	if err != nil {
		t.Fatalf("SetItemCategory() error = %v", err)
	}
	if retypedItem.Category != "work" {
		t.Fatalf("SetItemCategory() category = %q, want %q", retypedItem.Category, "work")
	}

	filteredHistory, err := dbStore.ListTextHistory(ctx, "work")
	if err != nil {
		t.Fatalf("ListTextHistory(work) error = %v", err)
	}
	if len(filteredHistory) != 1 || filteredHistory[0].ID != secondItem.ID {
		t.Fatalf("ListTextHistory(work) = %+v, want only second item", filteredHistory)
	}

	if err := dbStore.DeleteTextItem(ctx, firstItem.ID); err != nil {
		t.Fatalf("DeleteTextItem() error = %v", err)
	}

	_, err = dbStore.GetTextItemByID(ctx, firstItem.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetTextItemByID() after delete error = %v, want ErrNotFound", err)
	}

	history, err = dbStore.ListTextHistory(ctx, "")
	if err != nil {
		t.Fatalf("ListTextHistory() after delete error = %v", err)
	}
	if len(history) != 1 || history[0].ID != secondItem.ID {
		t.Fatalf("ListTextHistory() after delete = %+v, want only second item", history)
	}
}
