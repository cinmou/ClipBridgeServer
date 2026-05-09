// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cinmou/ClipBridgeServer/internal/store"
)

type textClipboardRequest struct {
	Text             string `json:"text"`
	Content          string `json:"content"`
	SourceDeviceID   string `json:"source_device_id"`
	SourceDeviceName string `json:"source_device_name"`
}

type categoryRequest struct {
	Name string `json:"name"`
}

type itemCategoryRequest struct {
	Category string `json:"category"`
}

type clipboardItemResponse struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	Text       string `json:"text"`
	IsFavorite bool   `json:"is_favorite"`
	Category   string `json:"category"`
	SourceID   string `json:"source_device_id,omitempty"`
	SourceName string `json:"source_device_name,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type clipboardHistoryResponse struct {
	Items []clipboardItemResponse `json:"items"`
}

type categoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type categoriesResponse struct {
	Items []categoryResponse `json:"items"`
}

func (r *Router) handleClipboardText(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	var payload textClipboardRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		statusCode, message := normalizeDecodeError(err)
		writeJSONError(w, statusCode, message)
		return
	}

	if err := ensureRequestFullyConsumed(req.Body); err != nil {
		statusCode, message := normalizeDecodeError(err)
		writeJSONError(w, statusCode, message)
		return
	}

	if strings.TrimSpace(payload.Text) == "" {
		payload.Text = payload.Content
	}

	if err := r.validateText(payload.Text); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := store.CreateTextItemInput{
		Text:             payload.Text,
		SourceDeviceID:   payload.SourceDeviceID,
		SourceDeviceName: payload.SourceDeviceName,
	}
	if r.cleaner != nil {
		settings, err := r.cleaner.GetSettings(req.Context())
		if err == nil {
			input.ExpiresAt = time.Now().UTC().Add(time.Duration(settings.TTLHours) * time.Hour).Format(time.RFC3339)
		}
	}

	item, err := r.store.CreateTextItem(req.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "create text clipboard item failed")
		return
	}

	writeJSONData(w, http.StatusCreated, toClipboardItemResponse(item))
}

func (r *Router) handleClipboardLatest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	item, err := r.store.GetLatestTextItem(req.Context())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "no text clipboard item found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "load latest text clipboard item failed")
		return
	}

	writeJSONData(w, http.StatusOK, toClipboardItemResponse(item))
}

func (r *Router) handleClipboardHistory(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	categoryName := strings.TrimSpace(req.URL.Query().Get("category"))
	items, err := r.store.ListTextHistory(req.Context(), categoryName)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusBadRequest, "category not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "load text clipboard history failed")
		return
	}

	writeJSONData(w, http.StatusOK, clipboardHistoryResponse{
		Items: toClipboardItemResponses(items),
	})
}

// handleClipboardItemRoutes keeps item-id endpoints and item sub-actions under
// one route prefix. That makes it easier to grow the API without exploding the
// router table every time one item-specific action is added.
func (r *Router) handleClipboardItemRoutes(w http.ResponseWriter, req *http.Request) {
	id, action, err := parseClipboardItemRoute(req.URL.Path)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	switch action {
	case "":
		r.handleClipboardItemByID(w, req, id)
	case "favorite":
		r.handleClipboardItemFavorite(w, req, id)
	case "category":
		r.handleClipboardItemCategory(w, req, id)
	default:
		writeJSONError(w, http.StatusNotFound, "clipboard item not found")
	}
}

func (r *Router) handleClipboardItemByID(w http.ResponseWriter, req *http.Request, id int64) {
	switch req.Method {
	case http.MethodGet:
		item, err := r.store.GetTextItemByID(req.Context(), id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, "text clipboard item not found")
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "load text clipboard item failed")
			return
		}

		writeJSONData(w, http.StatusOK, toClipboardItemResponse(item))
	case http.MethodDelete:
		if err := r.store.DeleteTextItem(req.Context(), id); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, "text clipboard item not found")
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "delete text clipboard item failed")
			return
		}

		writeJSONData(w, http.StatusOK, map[string]any{
			"deleted": true,
			"id":      id,
		})
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodDelete)
	}
}

func (r *Router) handleClipboardItemFavorite(w http.ResponseWriter, req *http.Request, id int64) {
	var (
		favorite bool
		methods  = []string{http.MethodPost, http.MethodDelete}
	)

	switch req.Method {
	case http.MethodPost:
		favorite = true
	case http.MethodDelete:
		favorite = false
	default:
		writeMethodNotAllowed(w, methods...)
		return
	}

	item, err := r.store.SetFavorite(req.Context(), id, favorite)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "text clipboard item not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "update favorite state failed")
		return
	}

	writeJSONData(w, http.StatusOK, toClipboardItemResponse(item))
}

func (r *Router) handleFavorites(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	items, err := r.store.ListFavorites(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "load favorites failed")
		return
	}

	writeJSONData(w, http.StatusOK, clipboardHistoryResponse{
		Items: toClipboardItemResponses(items),
	})
}

func (r *Router) handleCategories(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.handleListCategories(w, req)
	case http.MethodPost:
		r.handleCreateCategory(w, req)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (r *Router) handleListCategories(w http.ResponseWriter, req *http.Request) {
	categories, err := r.store.ListCategories(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "load categories failed")
		return
	}

	response := categoriesResponse{
		Items: make([]categoryResponse, 0, len(categories)),
	}
	for _, category := range categories {
		response.Items = append(response.Items, toCategoryResponse(category))
	}

	writeJSONData(w, http.StatusOK, response)
}

func (r *Router) handleCreateCategory(w http.ResponseWriter, req *http.Request) {
	var payload categoryRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		statusCode, message := normalizeDecodeError(err)
		writeJSONError(w, statusCode, message)
		return
	}

	if err := ensureRequestFullyConsumed(req.Body); err != nil {
		statusCode, message := normalizeDecodeError(err)
		writeJSONError(w, statusCode, message)
		return
	}

	categoryName := strings.TrimSpace(payload.Name)
	if categoryName == "" {
		writeJSONError(w, http.StatusBadRequest, "name must not be empty")
		return
	}

	category, err := r.store.CreateCategory(req.Context(), categoryName)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			writeJSONError(w, http.StatusConflict, "category already exists")
		default:
			writeJSONError(w, http.StatusInternalServerError, "create category failed")
		}
		return
	}

	writeJSONData(w, http.StatusCreated, toCategoryResponse(*category))
}

func (r *Router) handleClipboardItemCategory(w http.ResponseWriter, req *http.Request, id int64) {
	if req.Method != http.MethodPatch {
		writeMethodNotAllowed(w, http.MethodPatch)
		return
	}

	var payload itemCategoryRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		statusCode, message := normalizeDecodeError(err)
		writeJSONError(w, statusCode, message)
		return
	}

	if err := ensureRequestFullyConsumed(req.Body); err != nil {
		statusCode, message := normalizeDecodeError(err)
		writeJSONError(w, statusCode, message)
		return
	}

	categoryName := strings.TrimSpace(payload.Category)
	if categoryName == "" {
		writeJSONError(w, http.StatusBadRequest, "category must not be empty")
		return
	}

	item, err := r.store.SetItemCategory(req.Context(), id, categoryName)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound) && strings.Contains(err.Error(), "category"):
			writeJSONError(w, http.StatusBadRequest, "category not found")
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, "text clipboard item not found")
		default:
			writeJSONError(w, http.StatusInternalServerError, "update item category failed")
		}
		return
	}

	writeJSONData(w, http.StatusOK, toClipboardItemResponse(item))
}

func (r *Router) validateText(text string) error {
	if !utf8.ValidString(text) {
		return fmt.Errorf("text must be valid UTF-8")
	}

	textBytes := len([]byte(text))
	if textBytes < r.config.Limits.MinTextBytes {
		return fmt.Errorf("text must be at least %d bytes", r.config.Limits.MinTextBytes)
	}

	if textBytes > r.config.Limits.MaxTextBytes {
		return fmt.Errorf("text must be at most %d bytes", r.config.Limits.MaxTextBytes)
	}

	return nil
}

func parseClipboardItemRoute(path string) (int64, string, error) {
	const prefix = "/api/clipboard/items/"

	if !strings.HasPrefix(path, prefix) {
		return 0, "", fmt.Errorf("clipboard item not found")
	}

	trimmed := strings.TrimPrefix(path, prefix)
	if trimmed == "" {
		return 0, "", fmt.Errorf("clipboard item not found")
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) > 2 {
		return 0, "", fmt.Errorf("clipboard item not found")
	}

	id, err := strconv.ParseInt(segments[0], 10, 64)
	if err != nil || id <= 0 {
		return 0, "", fmt.Errorf("clipboard item not found")
	}

	if len(segments) == 1 {
		return id, "", nil
	}
	if segments[1] == "" {
		return 0, "", fmt.Errorf("clipboard item not found")
	}

	return id, segments[1], nil
}

func toClipboardItemResponses(items []store.ClipboardItem) []clipboardItemResponse {
	response := make([]clipboardItemResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toClipboardItemResponse(&item))
	}
	return response
}

func toClipboardItemResponse(item *store.ClipboardItem) clipboardItemResponse {
	return clipboardItemResponse{
		ID:         item.ID,
		Type:       item.ItemType,
		Text:       item.TextContent,
		IsFavorite: item.IsFavorite,
		Category:   item.Category,
		SourceID:   item.SourceDeviceID,
		SourceName: item.SourceDeviceName,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}
}

func toCategoryResponse(category store.Category) categoryResponse {
	return categoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func ensureRequestFullyConsumed(body io.Reader) error {
	var extra json.RawMessage
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("request body must contain exactly one JSON object")
		}
		return err
	}

	return nil
}

func normalizeDecodeError(err error) (int, string) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return http.StatusRequestEntityTooLarge, "request body is too large"
	}

	switch {
	case errors.Is(err, io.EOF):
		return http.StatusBadRequest, "request body must not be empty"
	default:
		return http.StatusBadRequest, "request body must be valid JSON"
	}
}
