// SPDX-License-Identifier: GPL-3.0-only

package webdav

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/store"
)

const syncVersion = 1

// Service owns the manual WebDAV sync preview. It intentionally avoids a
// background worker in this phase so behavior stays explicit and easy to debug.
type Service struct {
	store      *store.SQLiteStore
	dataDir    string
	httpClient *http.Client
}

type manifest struct {
	Version     int                    `json:"version"`
	UpdatedAt   string                 `json:"updated_at"`
	LastSyncAt  string                 `json:"last_sync_at"`
	ItemCount   int                    `json:"item_count"`
	Items       map[string]manifestRef `json:"items"`
}

type manifestRef struct {
	Type      string `json:"type"`
	UpdatedAt string `json:"updated_at"`
	FileName  string `json:"file_name,omitempty"`
}

type syncItem struct {
	Key              string `json:"key"`
	Type             string `json:"type"`
	Text             string `json:"text,omitempty"`
	URL              string `json:"url,omitempty"`
	Category         string `json:"category"`
	IsFavorite       bool   `json:"is_favorite"`
	SourceDeviceID   string `json:"source_device_id,omitempty"`
	SourceDeviceName string `json:"source_device_name,omitempty"`
	Filename         string `json:"filename,omitempty"`
	MIMEType         string `json:"mime_type,omitempty"`
	SHA256           string `json:"sha256,omitempty"`
	SizeBytes        int64  `json:"size_bytes"`
	ExpiresAt        string `json:"expires_at,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// NewService builds the manual WebDAV sync service.
func NewService(dbStore *store.SQLiteStore, cfg *config.Config) *Service {
	return &Service{
		store:   dbStore,
		dataDir: cfg.Storage.DataDir,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetHTTPClient lets tests replace the network transport while production keeps
// using a normal timeout-bound client.
func (s *Service) SetHTTPClient(client *http.Client) {
	if s == nil || client == nil {
		return
	}
	s.httpClient = client
}

func (s *Service) GetSettings(ctx context.Context) (store.WebDAVSettings, error) {
	return s.store.LoadWebDAVSettings(ctx)
}

func (s *Service) UpdateSettings(ctx context.Context, settings store.WebDAVSettings) (store.WebDAVSettings, error) {
	settings.URL = strings.TrimSpace(settings.URL)
	settings.Username = strings.TrimSpace(settings.Username)
	settings.BasePath = normalizeBasePath(settings.BasePath)
	if err := settings.Validate(); err != nil {
		return store.WebDAVSettings{}, err
	}
	if err := s.store.SaveWebDAVSettings(ctx, settings); err != nil {
		return store.WebDAVSettings{}, err
	}
	return settings, nil
}

func (s *Service) GetStatus(ctx context.Context) (store.WebDAVSyncStatus, error) {
	return s.store.LoadWebDAVSyncStatus(ctx)
}

func (s *Service) TestConnection(ctx context.Context) (store.WebDAVSyncStatus, error) {
	settings, err := s.GetSettings(ctx)
	if err != nil {
		return store.WebDAVSyncStatus{}, err
	}
	if err := settings.Validate(); err != nil {
		return store.WebDAVSyncStatus{}, err
	}

	client, err := newClient(settings, s.httpClient)
	if err != nil {
		return store.WebDAVSyncStatus{}, err
	}

	status, _ := s.GetStatus(ctx)
	status.TestedAt = time.Now().UTC().Format(time.RFC3339)
	if err := client.ping(ctx); err != nil {
		status.LastTestSuccess = false
		status.LastTestError = err.Error()
		if saveErr := s.store.SaveWebDAVSyncStatus(ctx, status); saveErr != nil {
			return status, saveErr
		}
		return status, err
	}

	status.LastTestSuccess = true
	status.LastTestError = ""
	if err := s.store.SaveWebDAVSyncStatus(ctx, status); err != nil {
		return status, err
	}
	return status, nil
}

func (s *Service) RunSync(ctx context.Context) (store.WebDAVSyncStatus, error) {
	settings, err := s.GetSettings(ctx)
	if err != nil {
		return store.WebDAVSyncStatus{}, err
	}
	if !settings.Enabled {
		return store.WebDAVSyncStatus{}, fmt.Errorf("webdav sync is disabled")
	}
	if err := settings.Validate(); err != nil {
		return store.WebDAVSyncStatus{}, err
	}

	client, err := newClient(settings, s.httpClient)
	if err != nil {
		return store.WebDAVSyncStatus{}, err
	}
	if err := client.ensureCollections(ctx); err != nil {
		return s.saveSyncError(ctx, err)
	}

	items, err := s.store.ListClipboardHistory(ctx, "")
	if err != nil {
		return s.saveSyncError(ctx, err)
	}

	localByKey := make(map[string]store.ClipboardItem, len(items))
	exported := make(map[string]syncItem, len(items))
	for _, item := range items {
		key := buildSyncKey(item)
		localByKey[key] = item
		exported[key] = toSyncItem(item, key)
	}

	remoteManifest, _ := client.loadManifest(ctx)
	if remoteManifest.Items == nil {
		remoteManifest.Items = make(map[string]manifestRef)
	}

	status, _ := s.GetStatus(ctx)
	status.LastSyncAt = time.Now().UTC().Format(time.RFC3339)
	status.LastError = ""
	status.LastMessage = ""
	status.PushedItems = 0
	status.PulledItems = 0
	status.PushedFiles = 0
	status.PulledFiles = 0
	status.ConflictSkips = 0
	status.LocalItemCount = len(items)

	for key, item := range exported {
		if err := client.putJSON(ctx, path.Join("items", key+".json"), item); err != nil {
			return s.saveSyncError(ctx, err)
		}
		status.PushedItems++

		if item.Type == "image" || item.Type == "file" {
			localItem := localByKey[key]
			if localItem.LocalPath != "" {
				if err := client.putFile(ctx, path.Join("files", key+".bin"), localItem.LocalPath, chooseRemoteContentType(item.MIMEType)); err != nil {
					return s.saveSyncError(ctx, err)
				}
				status.PushedFiles++
			}
		}

		remoteManifest.Items[key] = manifestRef{
			Type:      item.Type,
			UpdatedAt: item.UpdatedAt,
			FileName:  item.Filename,
		}
	}

	keys := make([]string, 0, len(remoteManifest.Items))
	for key := range remoteManifest.Items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if _, ok := localByKey[key]; ok {
			status.ConflictSkips++
			continue
		}
		remoteItem, err := client.loadItem(ctx, key)
		if err != nil {
			return s.saveSyncError(ctx, err)
		}
		if err := s.importRemoteItem(ctx, client, remoteItem); err != nil {
			return s.saveSyncError(ctx, err)
		}
		status.PulledItems++
		if remoteItem.Type == "image" || remoteItem.Type == "file" {
			status.PulledFiles++
		}
	}

	remoteManifest.Version = syncVersion
	remoteManifest.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	remoteManifest.LastSyncAt = remoteManifest.UpdatedAt
	remoteManifest.ItemCount = len(remoteManifest.Items)
	status.RemoteItemCount = remoteManifest.ItemCount
	status.LastSuccessAt = remoteManifest.UpdatedAt
	status.LastMessage = fmt.Sprintf("pushed %d items and pulled %d items", status.PushedItems, status.PulledItems)

	if err := client.putJSON(ctx, "manifest.json", remoteManifest); err != nil {
		return s.saveSyncError(ctx, err)
	}
	if err := s.store.SaveWebDAVSyncStatus(ctx, status); err != nil {
		return status, err
	}
	return status, nil
}

func (s *Service) saveSyncError(ctx context.Context, runErr error) (store.WebDAVSyncStatus, error) {
	status, _ := s.GetStatus(ctx)
	status.LastSyncAt = time.Now().UTC().Format(time.RFC3339)
	status.LastError = runErr.Error()
	status.LastMessage = ""
	_ = s.store.SaveWebDAVSyncStatus(ctx, status)
	return status, runErr
}

func (s *Service) importRemoteItem(ctx context.Context, client *client, item syncItem) error {
	if item.Key == "" {
		return fmt.Errorf("remote item sync key is missing")
	}

	input := store.ImportClipboardItemInput{
		ItemType:         item.Type,
		TextContent:      item.Text,
		Category:         item.Category,
		SourceDeviceID:   item.SourceDeviceID,
		SourceDeviceName: item.SourceDeviceName,
		Filename:         item.Filename,
		MIMEType:         item.MIMEType,
		SHA256:           item.SHA256,
		URL:              item.URL,
		SizeBytes:        item.SizeBytes,
		ExpiresAt:        item.ExpiresAt,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
		IsFavorite:       item.IsFavorite,
		SyncKey:          item.Key,
	}

	if item.Type == "image" || item.Type == "file" {
		localPath, err := client.downloadFile(ctx, path.Join("files", item.Key+".bin"), s.dataDir, item.Key, item.Filename)
		if err != nil {
			return err
		}
		input.LocalPath = localPath
	}

	_, err := s.store.ImportClipboardItem(ctx, input)
	return err
}

func toSyncItem(item store.ClipboardItem, key string) syncItem {
	return syncItem{
		Key:              key,
		Type:             item.ItemType,
		Text:             item.TextContent,
		URL:              item.TextContent,
		Category:         item.Category,
		IsFavorite:       item.IsFavorite,
		SourceDeviceID:   item.SourceDeviceID,
		SourceDeviceName: item.SourceDeviceName,
		Filename:         item.Filename,
		MIMEType:         item.MIMEType,
		SHA256:           item.SHA256,
		SizeBytes:        item.SizeBytes,
		ExpiresAt:        item.ExpiresAt,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}
}

func buildSyncKey(item store.ClipboardItem) string {
	if strings.TrimSpace(item.SyncKey) != "" {
		return strings.TrimSpace(item.SyncKey)
	}

	input := strings.Join([]string{
		item.ItemType,
		item.TextContent,
		item.Category,
		item.SourceDeviceID,
		item.SourceDeviceName,
		item.Filename,
		item.MIMEType,
		item.SHA256,
		fmt.Sprintf("%d", item.SizeBytes),
		item.CreatedAt,
	}, "\n")
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

func normalizeBasePath(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return "ClipBridgeServer"
	}
	return raw
}

func chooseRemoteContentType(mimeType string) string {
	if strings.TrimSpace(mimeType) == "" {
		return "application/octet-stream"
	}
	return mimeType
}

type client struct {
	baseURL    *url.URL
	username   string
	password   string
	httpClient *http.Client
}

func newClient(settings store.WebDAVSettings, httpClient *http.Client) (*client, error) {
	parsedURL, err := url.Parse(settings.URL)
	if err != nil {
		return nil, fmt.Errorf("parse webdav url: %w", err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("webdav url must be absolute")
	}
	basePath := normalizeBasePath(settings.BasePath)
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	if basePath != "" {
		parsedURL.Path = parsedURL.Path + "/" + basePath
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &client{
		baseURL:    parsedURL,
		username:   settings.Username,
		password:   settings.Password,
		httpClient: httpClient,
	}, nil
}

func (c *client) ping(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodOptions, "", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("test webdav connection: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("test webdav connection failed: %s", resp.Status)
	}
	return nil
}

func (c *client) ensureCollections(ctx context.Context) error {
	for _, collection := range []string{"", "items", "files"} {
		if err := c.mkcol(ctx, collection); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) mkcol(ctx context.Context, remotePath string) error {
	req, err := c.newRequest(ctx, "MKCOL", remotePath, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create remote collection %q: %w", remotePath, err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusCreated, http.StatusMethodNotAllowed, http.StatusConflict, http.StatusOK:
		return nil
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("create remote collection %q failed: %s %s", remotePath, resp.Status, strings.TrimSpace(string(body)))
	}
}

func (c *client) loadManifest(ctx context.Context) (manifest, error) {
	var result manifest
	err := c.getJSON(ctx, "manifest.json", &result)
	if errors.Is(err, os.ErrNotExist) {
		return manifest{Version: syncVersion, Items: map[string]manifestRef{}}, nil
	}
	return result, err
}

func (c *client) loadItem(ctx context.Context, key string) (syncItem, error) {
	var item syncItem
	if err := c.getJSON(ctx, path.Join("items", key+".json"), &item); err != nil {
		return syncItem{}, err
	}
	return item, nil
}

func (c *client) getJSON(ctx context.Context, remotePath string, target any) error {
	req, err := c.newRequest(ctx, http.MethodGet, remotePath, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("get remote %q: %w", remotePath, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return os.ErrNotExist
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("get remote %q failed: %s %s", remotePath, resp.Status, strings.TrimSpace(string(body)))
	}
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode remote %q: %w", remotePath, err)
	}
	return nil
}

func (c *client) putJSON(ctx context.Context, remotePath string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode remote json %q: %w", remotePath, err)
	}
	req, err := c.newRequest(ctx, http.MethodPut, remotePath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("put remote %q: %w", remotePath, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("put remote %q failed: %s %s", remotePath, resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func (c *client) putFile(ctx context.Context, remotePath, localPath, mimeType string) error {
	fileHandle, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local file %q: %w", localPath, err)
	}
	defer fileHandle.Close()

	req, err := c.newRequest(ctx, http.MethodPut, remotePath, fileHandle)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", chooseRemoteContentType(mimeType))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("put remote file %q: %w", remotePath, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("put remote file %q failed: %s %s", remotePath, resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func (c *client) downloadFile(ctx context.Context, remotePath, dataDir, key, filename string) (string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, remotePath, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download remote file %q: %w", remotePath, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("download remote file %q failed: %s %s", remotePath, resp.Status, strings.TrimSpace(string(body)))
	}

	dir := filepath.Join(dataDir, "uploads", "webdav")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create webdav download directory: %w", err)
	}

	safeName := sanitizeFilename(filename)
	if safeName == "" {
		safeName = key + ".bin"
	}
	localPath := filepath.Join(dir, key+"-"+safeName)
	fileHandle, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("create downloaded file %q: %w", localPath, err)
	}
	defer fileHandle.Close()

	if _, err := io.Copy(fileHandle, resp.Body); err != nil {
		return "", fmt.Errorf("write downloaded file %q: %w", localPath, err)
	}
	return localPath, nil
}

func (c *client) newRequest(ctx context.Context, method, remotePath string, body io.Reader) (*http.Request, error) {
	targetURL := *c.baseURL
	if remotePath != "" {
		targetURL.Path = strings.TrimSuffix(c.baseURL.Path, "/") + "/" + strings.TrimPrefix(remotePath, "/")
	}
	req, err := http.NewRequestWithContext(ctx, method, targetURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create %s request for %q: %w", method, remotePath, err)
	}
	req.SetBasicAuth(c.username, c.password)
	return req, nil
}

func sanitizeFilename(filename string) string {
	filename = strings.TrimSpace(filename)
	filename = strings.ReplaceAll(filename, "/", "-")
	filename = strings.ReplaceAll(filename, "\\", "-")
	return filename
}
