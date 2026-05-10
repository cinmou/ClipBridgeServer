// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/cleanup"
	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/store"
	"github.com/cinmou/ClipBridgeServer/internal/webdav"
	webui "github.com/cinmou/ClipBridgeServer/web"
)

func TestHealthDoesNotRequireToken(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	response := performRequest(t, testContext.router, http.MethodGet, "/api/health", nil, "")

	if response.Code != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want %d", response.Code, http.StatusOK)
	}

	var payload struct {
		Data healthResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode health response error = %v", err)
	}
	if !payload.Data.OK || payload.Data.Version != version {
		t.Fatalf("health payload = %+v, want ok=true version=%q", payload.Data, version)
	}
}

func TestRootServesEmbeddedWebUI(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	response := performRequest(t, testContext.router, http.MethodGet, "/", nil, "")

	if response.Code != http.StatusOK {
		t.Fatalf("GET / status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "ClipBridge Console") {
		t.Fatalf("GET / body = %q, want embedded web ui title", response.Body.String())
	}
}

func TestPairingFlowAndDeviceTokenClipboardAccess(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)

	pairingCodeResult := performRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/auth/pairing-codes",
		nil,
		adminToken,
	)
	if pairingCodeResult.Code != http.StatusCreated {
		t.Fatalf("POST /api/auth/pairing-codes status = %d, want %d", pairingCodeResult.Code, http.StatusCreated)
	}

	var pairingPayload struct {
		Data pairingCodeResponse `json:"data"`
	}
	if err := json.Unmarshal(pairingCodeResult.Body.Bytes(), &pairingPayload); err != nil {
		t.Fatalf("decode pairing code response error = %v", err)
	}
	if pairingPayload.Data.PairingCode == "" {
		t.Fatalf("pairing code should not be empty")
	}

	pairResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/auth/pair",
		map[string]string{
			"pairing_code": pairingPayload.Data.PairingCode,
			"device_name":  "Test Mac",
		},
		"",
	)
	if pairResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/auth/pair status = %d, want %d", pairResponse.Code, http.StatusCreated)
	}

	var pairedPayload struct {
		Data pairDeviceResponse `json:"data"`
	}
	if err := json.Unmarshal(pairResponse.Body.Bytes(), &pairedPayload); err != nil {
		t.Fatalf("decode pair response error = %v", err)
	}
	if pairedPayload.Data.DeviceToken == "" {
		t.Fatalf("device token should not be empty")
	}

	createResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{"text": "hello"},
		pairedPayload.Data.DeviceToken,
	)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/text with device token status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	createdItem := decodeClipboardItemData(t, createResponse)
	if createdItem.Category != "text" {
		t.Fatalf("created item category = %q, want %q", createdItem.Category, "text")
	}
	if createdItem.SourceName != "" {
		t.Fatalf("created item source name = %q, want empty", createdItem.SourceName)
	}

	latestResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/latest", nil, pairedPayload.Data.DeviceToken)
	if latestResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/clipboard/latest with device token status = %d, want %d", latestResponse.Code, http.StatusOK)
	}

	itemPath := "/api/clipboard/items/" + jsonNumber(createdItem.ID)
	deleteResponse := performRequest(t, testContext.router, http.MethodDelete, itemPath, nil, pairedPayload.Data.DeviceToken)
	if deleteResponse.Code != http.StatusOK {
		t.Fatalf("DELETE %s with device token status = %d, want %d", itemPath, deleteResponse.Code, http.StatusOK)
	}
}

func TestQuickClipboardUploadPayloadUsesExistingClipboardAPI(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)

	createResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{
			"content":            "web txt",
			"source_device_id":   "web-ui",
			"source_device_name": "Web UI",
		},
		adminToken,
	)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/text quick payload status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	item := decodeClipboardItemData(t, createResponse)
	if item.Text != "web txt" {
		t.Fatalf("created item text = %q, want %q", item.Text, "web txt")
	}
	if item.SourceID != "web-ui" || item.SourceName != "Web UI" {
		t.Fatalf("created item source = (%q, %q), want (%q, %q)", item.SourceID, item.SourceName, "web-ui", "Web UI")
	}
}

func TestLinkAndFileClipboardFlow(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)

	linkResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/link",
		map[string]string{
			"url":                "https://example.com/demo",
			"source_device_id":   "web-ui",
			"source_device_name": "Web UI",
		},
		adminToken,
	)
	if linkResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/link status = %d, want %d", linkResponse.Code, http.StatusCreated)
	}
	linkItem := decodeClipboardItemData(t, linkResponse)
	if linkItem.Type != "link" || linkItem.Category != "link" {
		t.Fatalf("link item = %+v, want type/category link", linkItem)
	}

	imageBytes := []byte{
		0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n',
		0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	}
	fileResponse := performMultipartUpload(
		t,
		testContext.router,
		"/api/clipboard/file",
		"preview.png",
		imageBytes,
		map[string]string{
			"source_device_id":   "web-ui",
			"source_device_name": "Web UI",
		},
		adminToken,
	)
	if fileResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/file status = %d, want %d", fileResponse.Code, http.StatusCreated)
	}
	fileItem := decodeClipboardItemData(t, fileResponse)
	if fileItem.Type != "image" {
		t.Fatalf("file item type = %q, want %q", fileItem.Type, "image")
	}
	if fileItem.DownloadURL == "" || fileItem.PreviewURL == "" {
		t.Fatalf("file item urls = %+v, want download and preview urls", fileItem)
	}

	historyResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/history", nil, adminToken)
	if historyResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/clipboard/history status = %d, want %d", historyResponse.Code, http.StatusOK)
	}

	var historyPayload struct {
		Data clipboardHistoryResponse `json:"data"`
	}
	if err := json.Unmarshal(historyResponse.Body.Bytes(), &historyPayload); err != nil {
		t.Fatalf("decode history response error = %v", err)
	}
	if len(historyPayload.Data.Items) != 2 {
		t.Fatalf("history items len = %d, want 2", len(historyPayload.Data.Items))
	}

	downloadResponse := performRequest(t, testContext.router, http.MethodGet, fileItem.DownloadURL, nil, adminToken)
	if downloadResponse.Code != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d", fileItem.DownloadURL, downloadResponse.Code, http.StatusOK)
	}
	if got := downloadResponse.Header().Get("Content-Type"); got == "" || got == "application/octet-stream" {
		t.Fatalf("download content-type = %q, want detected image mime", got)
	}
	if !bytes.Equal(downloadResponse.Body.Bytes(), imageBytes) {
		t.Fatalf("downloaded bytes = %v, want %v", downloadResponse.Body.Bytes(), imageBytes)
	}
}

func TestSettingsLimitsAndAdminTokenCanBeUpdated(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)

	limitsResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPatch,
		"/api/settings/limits",
		store.LimitsSettings{
			MinTextBytes:    1,
			MaxTextBytes:    4,
			MinImageBytes:   1,
			MaxImageBytes:   64,
			MinFileBytes:    1,
			MaxFileBytes:    64,
			MinLinkBytes:    1,
			MaxLinkBytes:    32,
			MaxRequestBytes: 512,
		},
		adminToken,
	)
	if limitsResponse.Code != http.StatusOK {
		t.Fatalf("PATCH /api/settings/limits status = %d, want %d", limitsResponse.Code, http.StatusOK)
	}

	tooLargeTextResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{"text": "12345"},
		adminToken,
	)
	if tooLargeTextResponse.Code != http.StatusBadRequest {
		t.Fatalf("POST /api/clipboard/text after limits patch status = %d, want %d", tooLargeTextResponse.Code, http.StatusBadRequest)
	}

	settingsResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPatch,
		"/api/settings",
		map[string]string{"admin_token": "new-admin-token"},
		adminToken,
	)
	if settingsResponse.Code != http.StatusOK {
		t.Fatalf("PATCH /api/settings status = %d, want %d", settingsResponse.Code, http.StatusOK)
	}

	oldTokenResponse := performRequest(t, testContext.router, http.MethodGet, "/api/settings", nil, adminToken)
	if oldTokenResponse.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/settings with old token status = %d, want %d", oldTokenResponse.Code, http.StatusUnauthorized)
	}

	newTokenResponse := performRequest(t, testContext.router, http.MethodGet, "/api/settings", nil, "new-admin-token")
	if newTokenResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/settings with new token status = %d, want %d", newTokenResponse.Code, http.StatusOK)
	}
}

func TestFavoritesAndCategoriesFlow(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)

	createResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{"text": "fav me"},
		adminToken,
	)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/text status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	createdItem := decodeClipboardItemData(t, createResponse)
	if createdItem.Category != "text" {
		t.Fatalf("created item category = %q, want %q", createdItem.Category, "text")
	}
	if createdItem.IsFavorite {
		t.Fatalf("created item isFavorite = %v, want false", createdItem.IsFavorite)
	}

	favoritePath := "/api/clipboard/items/" + jsonNumber(createdItem.ID) + "/favorite"
	favoriteResponse := performRequest(t, testContext.router, http.MethodPost, favoritePath, nil, adminToken)
	if favoriteResponse.Code != http.StatusOK {
		t.Fatalf("POST %s status = %d, want %d", favoritePath, favoriteResponse.Code, http.StatusOK)
	}

	favoritedItem := decodeClipboardItemData(t, favoriteResponse)
	if !favoritedItem.IsFavorite {
		t.Fatalf("favorited item isFavorite = %v, want true", favoritedItem.IsFavorite)
	}

	favoritesResponse := performRequest(t, testContext.router, http.MethodGet, "/api/favorites", nil, adminToken)
	if favoritesResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/favorites status = %d, want %d", favoritesResponse.Code, http.StatusOK)
	}

	var favoritesPayload struct {
		Data clipboardHistoryResponse `json:"data"`
	}
	if err := json.Unmarshal(favoritesResponse.Body.Bytes(), &favoritesPayload); err != nil {
		t.Fatalf("decode favorites response error = %v", err)
	}
	if len(favoritesPayload.Data.Items) != 1 || favoritesPayload.Data.Items[0].ID != createdItem.ID {
		t.Fatalf("favorites payload = %+v, want created item only", favoritesPayload.Data.Items)
	}

	createCategoryResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/categories",
		map[string]string{"name": "work"},
		adminToken,
	)
	if createCategoryResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/categories status = %d, want %d", createCategoryResponse.Code, http.StatusCreated)
	}

	categoryPatchPath := "/api/clipboard/items/" + jsonNumber(createdItem.ID) + "/category"
	categoryPatchResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPatch,
		categoryPatchPath,
		map[string]string{"category": "work"},
		adminToken,
	)
	if categoryPatchResponse.Code != http.StatusOK {
		t.Fatalf("PATCH %s status = %d, want %d", categoryPatchPath, categoryPatchResponse.Code, http.StatusOK)
	}

	updatedItem := decodeClipboardItemData(t, categoryPatchResponse)
	if updatedItem.Category != "work" {
		t.Fatalf("patched item category = %q, want %q", updatedItem.Category, "work")
	}

	filteredHistoryResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/history?category=work", nil, adminToken)
	if filteredHistoryResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/clipboard/history?category=work status = %d, want %d", filteredHistoryResponse.Code, http.StatusOK)
	}

	var filteredHistoryPayload struct {
		Data clipboardHistoryResponse `json:"data"`
	}
	if err := json.Unmarshal(filteredHistoryResponse.Body.Bytes(), &filteredHistoryPayload); err != nil {
		t.Fatalf("decode filtered history response error = %v", err)
	}
	if len(filteredHistoryPayload.Data.Items) != 1 || filteredHistoryPayload.Data.Items[0].Category != "work" {
		t.Fatalf("filtered history payload = %+v, want one work item", filteredHistoryPayload.Data.Items)
	}

	unfavoriteResponse := performRequest(t, testContext.router, http.MethodDelete, favoritePath, nil, adminToken)
	if unfavoriteResponse.Code != http.StatusOK {
		t.Fatalf("DELETE %s status = %d, want %d", favoritePath, unfavoriteResponse.Code, http.StatusOK)
	}

	unfavoritedItem := decodeClipboardItemData(t, unfavoriteResponse)
	if unfavoritedItem.IsFavorite {
		t.Fatalf("unfavorited item isFavorite = %v, want false", unfavoritedItem.IsFavorite)
	}
}

func TestPairingCodeCanOnlyBeUsedOnce(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	pairingCode := createPairingCodeForTest(t, testContext.router)

	firstPairResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/auth/pair",
		map[string]string{"pairing_code": pairingCode, "device_name": "Device One"},
		"",
	)
	if firstPairResponse.Code != http.StatusCreated {
		t.Fatalf("first POST /api/auth/pair status = %d, want %d", firstPairResponse.Code, http.StatusCreated)
	}

	secondPairResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/auth/pair",
		map[string]string{"pairing_code": pairingCode, "device_name": "Device Two"},
		"",
	)
	if secondPairResponse.Code != http.StatusBadRequest {
		t.Fatalf("second POST /api/auth/pair status = %d, want %d", secondPairResponse.Code, http.StatusBadRequest)
	}
	assertErrorMessage(t, secondPairResponse, "pairing code has already been used")
}

func TestExpiredPairingCodeFails(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	expiredCode, err := testContext.store.CreatePairingCode(testContext.ctx, time.Now().UTC().Add(-1*time.Minute))
	if err != nil {
		t.Fatalf("CreatePairingCode(expired) error = %v", err)
	}

	response := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/auth/pair",
		map[string]string{"pairing_code": expiredCode.Code, "device_name": "Late Device"},
		"",
	)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expired POST /api/auth/pair status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	assertErrorMessage(t, response, "pairing code has expired")
}

func TestClipboardRequiresValidToken(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)

	noTokenResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/latest", nil, "")
	if noTokenResponse.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/clipboard/latest without token status = %d, want %d", noTokenResponse.Code, http.StatusUnauthorized)
	}
	assertErrorMessage(t, noTokenResponse, "unauthorized")

	wrongTokenResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/latest", nil, "wrong-token")
	if wrongTokenResponse.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/clipboard/latest with wrong token status = %d, want %d", wrongTokenResponse.Code, http.StatusUnauthorized)
	}
	assertErrorMessage(t, wrongTokenResponse, "unauthorized")
}

func TestDevicesListAndRevocation(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	pairingCode := createPairingCodeForTest(t, testContext.router)
	pairResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/auth/pair",
		map[string]string{"pairing_code": pairingCode, "device_name": "Revokable Device"},
		"",
	)
	if pairResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/auth/pair status = %d, want %d", pairResponse.Code, http.StatusCreated)
	}

	var pairedPayload struct {
		Data pairDeviceResponse `json:"data"`
	}
	if err := json.Unmarshal(pairResponse.Body.Bytes(), &pairedPayload); err != nil {
		t.Fatalf("decode pair response error = %v", err)
	}

	devicesResult := performRequest(t, testContext.router, http.MethodGet, "/api/auth/devices", nil, adminToken)
	if devicesResult.Code != http.StatusOK {
		t.Fatalf("GET /api/auth/devices status = %d, want %d", devicesResult.Code, http.StatusOK)
	}

	var devicesPayload struct {
		Data devicesResponse `json:"data"`
	}
	if err := json.Unmarshal(devicesResult.Body.Bytes(), &devicesPayload); err != nil {
		t.Fatalf("decode devices response error = %v", err)
	}
	if len(devicesPayload.Data.Items) != 1 {
		t.Fatalf("devices list len = %d, want 1", len(devicesPayload.Data.Items))
	}

	deviceID := pairedPayload.Data.Device.ID
	revokePath := "/api/auth/devices/" + jsonNumber(deviceID)
	revokeResponse := performRequest(t, testContext.router, http.MethodDelete, revokePath, nil, adminToken)
	if revokeResponse.Code != http.StatusOK {
		t.Fatalf("DELETE %s status = %d, want %d", revokePath, revokeResponse.Code, http.StatusOK)
	}

	clipboardResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/latest", nil, pairedPayload.Data.DeviceToken)
	if clipboardResponse.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/clipboard/latest after revoke status = %d, want %d", clipboardResponse.Code, http.StatusUnauthorized)
	}
	assertErrorMessage(t, clipboardResponse, "unauthorized")
}

func TestClipboardValidationAndRequestTooLarge(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)

	emptyResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{"text": ""},
		adminToken,
	)
	if emptyResponse.Code != http.StatusBadRequest {
		t.Fatalf("empty POST status = %d, want %d", emptyResponse.Code, http.StatusBadRequest)
	}
	assertErrorMessage(t, emptyResponse, "text must be at least 1 bytes")

	tooLargeTextResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{"text": "123456789"},
		adminToken,
	)
	if tooLargeTextResponse.Code != http.StatusBadRequest {
		t.Fatalf("too large text POST status = %d, want %d", tooLargeTextResponse.Code, http.StatusBadRequest)
	}
	assertErrorMessage(t, tooLargeTextResponse, "text must be at most 8 bytes")

	oversizedJSON := []byte(`{"text":"` + strings.Repeat("a", 5000) + `"}`)
	tooLargeRequestResponse := performRequest(t, testContext.router, http.MethodPost, "/api/clipboard/text", oversizedJSON, adminToken)
	if tooLargeRequestResponse.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized request status = %d, want %d", tooLargeRequestResponse.Code, http.StatusRequestEntityTooLarge)
	}
	assertErrorMessage(t, tooLargeRequestResponse, "request body is too large")
}

func TestClipboardPreflightCORS(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	request := httptest.NewRequest(http.MethodOptions, "/api/clipboard/text", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	request.Header.Set("Access-Control-Request-Method", "POST")

	recorder := httptest.NewRecorder()
	testContext.router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("OPTIONS /api/clipboard/text status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
}

func TestCategoriesListContainsBuiltins(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	response := performRequest(t, testContext.router, http.MethodGet, "/api/categories", nil, adminToken)
	if response.Code != http.StatusOK {
		t.Fatalf("GET /api/categories status = %d, want %d", response.Code, http.StatusOK)
	}

	var payload struct {
		Data categoriesResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode categories response error = %v", err)
	}

	categoryNames := make(map[string]bool, len(payload.Data.Items))
	for _, item := range payload.Data.Items {
		categoryNames[item.Name] = true
	}

	for _, builtin := range []string{"text", "image", "link", "file"} {
		if !categoryNames[builtin] {
			t.Fatalf("expected built-in category %q in %+v", builtin, payload.Data.Items)
		}
	}
}

func TestCleanupSettingsAndManualRun(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	createResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{"text": "expire"},
		adminToken,
	)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/text status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	settingsResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPatch,
		"/api/settings/cleanup",
		store.CleanupSettings{
			TTLHours:        1,
			MaxItems:        2,
			MaxTotalSizeMB:  1,
			IntervalMinutes: 5,
			Enabled:         true,
		},
		adminToken,
	)
	if settingsResponse.Code != http.StatusOK {
		t.Fatalf("PATCH /api/settings/cleanup status = %d, want %d", settingsResponse.Code, http.StatusOK)
	}

	db, err := sql.Open("sqlite", testContext.dbPath)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if _, err := db.ExecContext(
		testContext.ctx,
		`UPDATE clipboard_items SET created_at = ?, expires_at = ? WHERE id = 1`,
		"2000-01-01T00:00:00Z",
		"2000-01-01T01:00:00Z",
	); err != nil {
		t.Fatalf("seed expired item error = %v", err)
	}

	runResponse := performRequest(t, testContext.router, http.MethodPost, "/api/admin/cleanup/run", nil, adminToken)
	if runResponse.Code != http.StatusOK {
		t.Fatalf("POST /api/admin/cleanup/run status = %d, want %d", runResponse.Code, http.StatusOK)
	}

	latestResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/latest", nil, adminToken)
	if latestResponse.Code != http.StatusNotFound {
		t.Fatalf("GET /api/clipboard/latest after cleanup status = %d, want %d", latestResponse.Code, http.StatusNotFound)
	}

	statusResponse := performRequest(t, testContext.router, http.MethodGet, "/api/admin/cleanup/status", nil, adminToken)
	if statusResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/admin/cleanup/status status = %d, want %d", statusResponse.Code, http.StatusOK)
	}

	storageResponse := performRequest(t, testContext.router, http.MethodGet, "/api/admin/storage/status", nil, adminToken)
	if storageResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/admin/storage/status status = %d, want %d", storageResponse.Code, http.StatusOK)
	}
}

func TestWebDAVSettingsConnectionAndManualSync(t *testing.T) {
	t.Parallel()

	testContext := newTestContext(t)
	webdavServer := newFakeWebDAVServer(t, "/dav", "demo", "secret")
	testContext.webdav.SetHTTPClient(&http.Client{Transport: webdavServer})

	settingsResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPatch,
		"/api/settings/webdav",
		store.WebDAVSettings{
			Enabled:  true,
			URL:      "http://webdav.test/dav",
			Username: "demo",
			Password: "secret",
			BasePath: "ClipBridgeServer",
		},
		adminToken,
	)
	if settingsResponse.Code != http.StatusOK {
		t.Fatalf("PATCH /api/settings/webdav status = %d, want %d", settingsResponse.Code, http.StatusOK)
	}

	testResponse := performRequest(t, testContext.router, http.MethodPost, "/api/admin/webdav/test", nil, adminToken)
	if testResponse.Code != http.StatusOK {
		t.Fatalf("POST /api/admin/webdav/test status = %d, want %d", testResponse.Code, http.StatusOK)
	}

	createTextResponse := performJSONRequest(
		t,
		testContext.router,
		http.MethodPost,
		"/api/clipboard/text",
		map[string]string{"text": "local"},
		adminToken,
	)
	if createTextResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/text status = %d, want %d", createTextResponse.Code, http.StatusCreated)
	}

	fileBytes := []byte("local-file-bytes")
	createFileResponse := performMultipartUpload(
		t,
		testContext.router,
		"/api/clipboard/file",
		"demo.txt",
		fileBytes,
		map[string]string{},
		adminToken,
	)
	if createFileResponse.Code != http.StatusCreated {
		t.Fatalf("POST /api/clipboard/file status = %d, want %d", createFileResponse.Code, http.StatusCreated)
	}

	webdavServer.putJSON(t, "/dav/ClipBridgeServer/manifest.json", map[string]any{
		"version":    1,
		"updated_at": "2026-05-09T12:00:00Z",
		"last_sync_at": "2026-05-09T12:00:00Z",
		"item_count": 2,
		"items": map[string]any{
			"remote-text": map[string]any{
				"type":       "text",
				"updated_at": "2026-05-09T12:00:00Z",
			},
			"remote-file": map[string]any{
				"type":       "file",
				"updated_at": "2026-05-09T12:05:00Z",
				"file_name":  "remote.bin",
			},
		},
	})
	webdavServer.putJSON(t, "/dav/ClipBridgeServer/items/remote-text.json", map[string]any{
		"key":         "remote-text",
		"type":        "text",
		"text":        "remote text",
		"category":    "text",
		"size_bytes":  11,
		"created_at":  "2026-05-09T12:00:00Z",
		"updated_at":  "2026-05-09T12:00:00Z",
		"is_favorite": false,
	})
	webdavServer.putJSON(t, "/dav/ClipBridgeServer/items/remote-file.json", map[string]any{
		"key":         "remote-file",
		"type":        "file",
		"category":    "file",
		"filename":    "remote.bin",
		"mime_type":   "application/octet-stream",
		"sha256":      "abc123",
		"size_bytes":  12,
		"created_at":  "2026-05-09T12:05:00Z",
		"updated_at":  "2026-05-09T12:05:00Z",
		"is_favorite": false,
	})
	webdavServer.putBytes("/dav/ClipBridgeServer/files/remote-file.bin", []byte("remote bytes"))

	syncResponse := performRequest(t, testContext.router, http.MethodPost, "/api/admin/webdav/sync", nil, adminToken)
	if syncResponse.Code != http.StatusOK {
		t.Fatalf("POST /api/admin/webdav/sync status = %d, want %d body=%s", syncResponse.Code, http.StatusOK, syncResponse.Body.String())
	}

	var syncPayload struct {
		Data store.WebDAVSyncStatus `json:"data"`
	}
	if err := json.Unmarshal(syncResponse.Body.Bytes(), &syncPayload); err != nil {
		t.Fatalf("decode webdav sync response error = %v", err)
	}
	if syncPayload.Data.PulledItems != 2 {
		t.Fatalf("pulled items = %d, want 2", syncPayload.Data.PulledItems)
	}
	if syncPayload.Data.PushedItems < 2 {
		t.Fatalf("pushed items = %d, want at least 2", syncPayload.Data.PushedItems)
	}

	historyResponse := performRequest(t, testContext.router, http.MethodGet, "/api/clipboard/history", nil, adminToken)
	if historyResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/clipboard/history status = %d, want %d", historyResponse.Code, http.StatusOK)
	}

	var historyPayload struct {
		Data clipboardHistoryResponse `json:"data"`
	}
	if err := json.Unmarshal(historyResponse.Body.Bytes(), &historyPayload); err != nil {
		t.Fatalf("decode history response error = %v", err)
	}
	if len(historyPayload.Data.Items) != 4 {
		t.Fatalf("history items len = %d, want 4", len(historyPayload.Data.Items))
	}

	foundRemoteText := false
	foundRemoteFile := false
	for _, item := range historyPayload.Data.Items {
		if item.Text == "remote text" {
			foundRemoteText = true
		}
		if item.Filename == "remote.bin" && item.Type == "file" {
			foundRemoteFile = true
		}
	}
	if !foundRemoteText || !foundRemoteFile {
		t.Fatalf("expected remote text and remote file in history, got %+v", historyPayload.Data.Items)
	}

	statusResponse := performRequest(t, testContext.router, http.MethodGet, "/api/admin/webdav/status", nil, adminToken)
	if statusResponse.Code != http.StatusOK {
		t.Fatalf("GET /api/admin/webdav/status status = %d, want %d", statusResponse.Code, http.StatusOK)
	}
}

const adminToken = "test-admin-token-123"

type testContext struct {
	ctx    context.Context
	router http.Handler
	store  *store.SQLiteStore
	dbPath string
	webdav *webdav.Service
}

func newTestContext(t *testing.T) testContext {
	t.Helper()

	cfg := config.Default()
	cfg.Auth.Token = adminToken
	cfg.Limits.MinTextBytes = 1
	cfg.Limits.MaxTextBytes = 8
	cfg.Limits.MaxRequestBytes = 2048
	cfg.Storage.TTLHours = 24
	cfg.Storage.MaxItems = 10
	cfg.Storage.MaxTotalSizeMB = 1
	cfg.Cleaner.Enabled = false
	cfg.Cleaner.IntervalMinutes = 5

	dbPath := filepath.Join(t.TempDir(), "clipbridge.db")
	dbStore, err := store.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	t.Cleanup(func() {
		_ = dbStore.Close()
	})

	cleanupService, err := cleanup.NewService(dbStore, cfg)
	if err != nil {
		t.Fatalf("cleanup.NewService() error = %v", err)
	}
	t.Cleanup(func() {
		cleanupService.Close()
	})

	webdavService := webdav.NewService(dbStore, cfg)

	return testContext{
		ctx:    context.Background(),
		router: NewRouter(dbStore, cfg, cleanupService, webdavService, webui.Handler()),
		store:  dbStore,
		dbPath: dbPath,
		webdav: webdavService,
	}
}

type fakeWebDAVServer struct {
	t        *testing.T
	username string
	password string
	basePath string
	mu       sync.Mutex
	files    map[string][]byte
}

func newFakeWebDAVServer(t *testing.T, basePath, username, password string) *fakeWebDAVServer {
	t.Helper()

	fake := &fakeWebDAVServer{
		t:        t,
		username: username,
		password: password,
		basePath: basePath,
		files:    make(map[string][]byte),
	}
	return fake
}

func (f *fakeWebDAVServer) putJSON(t *testing.T, remotePath string, payload any) {
	t.Helper()

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	f.putBytes(remotePath, raw)
}

func (f *fakeWebDAVServer) putBytes(remotePath string, body []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.files[remotePath] = append([]byte(nil), body...)
}

func (f *fakeWebDAVServer) RoundTrip(req *http.Request) (*http.Response, error) {
	user, pass, ok := req.BasicAuth()
	if !ok || user != f.username || pass != f.password {
		return fakeHTTPResponse(req, http.StatusUnauthorized, nil, ""), nil
	}
	if !strings.HasPrefix(req.URL.Path, f.basePath) {
		return fakeHTTPResponse(req, http.StatusNotFound, nil, ""), nil
	}

	switch req.Method {
	case http.MethodOptions:
		return fakeHTTPResponse(req, http.StatusOK, nil, ""), nil
	case "MKCOL":
		return fakeHTTPResponse(req, http.StatusCreated, nil, ""), nil
	case http.MethodPut:
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		f.putBytes(req.URL.Path, body)
		return fakeHTTPResponse(req, http.StatusCreated, nil, ""), nil
	case http.MethodGet:
		f.mu.Lock()
		body, ok := f.files[req.URL.Path]
		f.mu.Unlock()
		if !ok {
			return fakeHTTPResponse(req, http.StatusNotFound, nil, ""), nil
		}
		contentType := ""
		if strings.HasSuffix(req.URL.Path, ".json") {
			contentType = "application/json"
		}
		return fakeHTTPResponse(req, http.StatusOK, body, contentType), nil
	default:
		return fakeHTTPResponse(req, http.StatusMethodNotAllowed, nil, ""), nil
	}
}

func fakeHTTPResponse(req *http.Request, statusCode int, body []byte, contentType string) *http.Response {
	header := make(http.Header)
	if contentType != "" {
		header.Set("Content-Type", contentType)
	}
	return &http.Response{
		StatusCode: statusCode,
		Status:     strconv.Itoa(statusCode) + " " + http.StatusText(statusCode),
		Header:     header,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Request:    req,
	}
}

func createPairingCodeForTest(t *testing.T, handler http.Handler) string {
	t.Helper()

	response := performRequest(t, handler, http.MethodPost, "/api/auth/pairing-codes", nil, adminToken)
	if response.Code != http.StatusCreated {
		t.Fatalf("POST /api/auth/pairing-codes status = %d, want %d", response.Code, http.StatusCreated)
	}

	var payload struct {
		Data pairingCodeResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode pairing code response error = %v", err)
	}

	return payload.Data.PairingCode
}

func performJSONRequest(t *testing.T, handler http.Handler, method, path string, payload any, token string) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	return performRequest(t, handler, method, path, body, token)
}

func performRequest(t *testing.T, handler http.Handler, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder
}

func performMultipartUpload(
	t *testing.T,
	handler http.Handler,
	path string,
	filename string,
	fileBytes []byte,
	fields map[string]string,
	token string,
) *httptest.ResponseRecorder {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("WriteField(%q) error = %v", key, err)
		}
	}

	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := fileWriter.Write(fileBytes); err != nil {
		t.Fatalf("fileWriter.Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body.Bytes()))
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func decodeClipboardItemData(t *testing.T, response *httptest.ResponseRecorder) clipboardItemResponse {
	t.Helper()

	var payload struct {
		Data clipboardItemResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode clipboard item response error = %v", err)
	}

	return payload.Data
}

func assertErrorMessage(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	var payload struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error response error = %v", err)
	}
	if payload.Error.Message != want {
		t.Fatalf("error message = %q, want %q", payload.Error.Message, want)
	}
}

func jsonNumber(value int64) string {
	return strconv.FormatInt(value, 10)
}
