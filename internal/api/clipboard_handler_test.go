// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/cleanup"
	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/store"
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

	oversizedJSON := []byte(`{"text":"` + strings.Repeat("a", 600) + `"}`)
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

const adminToken = "test-admin-token-123"

type testContext struct {
	ctx    context.Context
	router http.Handler
	store  *store.SQLiteStore
	dbPath string
}

func newTestContext(t *testing.T) testContext {
	t.Helper()

	cfg := config.Default()
	cfg.Auth.Token = adminToken
	cfg.Limits.MinTextBytes = 1
	cfg.Limits.MaxTextBytes = 8
	cfg.Limits.MaxRequestBytes = 512
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

	return testContext{
		ctx:    context.Background(),
		router: NewRouter(dbStore, cfg, cleanupService, webui.Handler()),
		store:  dbStore,
		dbPath: dbPath,
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
