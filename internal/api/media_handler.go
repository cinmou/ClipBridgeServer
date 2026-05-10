// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/store"
)

type linkClipboardRequest struct {
	URL              string `json:"url"`
	SourceDeviceID   string `json:"source_device_id"`
	SourceDeviceName string `json:"source_device_name"`
}

func (r *Router) handleClipboardFile(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	fileInput, err := r.readMultipartFile(req)
	if err != nil {
		statusCode, message := mapMultipartError(err)
		writeJSONError(w, statusCode, message)
		return
	}

	item, err := r.store.CreateFileItem(req.Context(), fileInput)
	if err != nil {
		_ = os.Remove(fileInput.LocalPath)
		writeJSONError(w, http.StatusInternalServerError, "create file clipboard item failed")
		return
	}

	writeJSONData(w, http.StatusCreated, toClipboardItemResponse(item))
}

func (r *Router) handleClipboardLink(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	var payload linkClipboardRequest
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

	rawURL := strings.TrimSpace(payload.URL)
	if err := r.validateLink(req.Context(), rawURL); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := store.CreateLinkItemInput{
		URL:              rawURL,
		SourceDeviceID:   payload.SourceDeviceID,
		SourceDeviceName: payload.SourceDeviceName,
	}
	if r.cleaner != nil {
		settings, err := r.cleaner.GetSettings(req.Context())
		if err == nil {
			input.ExpiresAt = time.Now().UTC().Add(time.Duration(settings.TTLHours) * time.Hour).Format(time.RFC3339)
		}
	}

	item, err := r.store.CreateLinkItem(req.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "create link clipboard item failed")
		return
	}

	writeJSONData(w, http.StatusCreated, toClipboardItemResponse(item))
}

func (r *Router) handleClipboardItemFile(w http.ResponseWriter, req *http.Request, id int64) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	item, err := r.store.GetClipboardItemByID(req.Context(), id)
	if err != nil {
		if err == store.ErrNotFound {
			writeJSONError(w, http.StatusNotFound, "clipboard item not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "load clipboard item failed")
		return
	}

	if item.ItemType != "image" && item.ItemType != "file" {
		writeJSONError(w, http.StatusBadRequest, "clipboard item has no downloadable file")
		return
	}

	if item.LocalPath == "" {
		writeJSONError(w, http.StatusInternalServerError, "clipboard file path is missing")
		return
	}

	fileHandle, err := os.Open(item.LocalPath)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "open clipboard file failed")
		return
	}
	defer fileHandle.Close()

	disposition := "attachment"
	if item.ItemType == "image" {
		disposition = "inline"
	}

	filename := item.Filename
	if filename == "" {
		filename = filepath.Base(item.LocalPath)
	}

	w.Header().Set("Content-Type", chooseContentType(item.MIMEType))
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=%q", disposition, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", item.SizeBytes))
	_, _ = io.Copy(w, fileHandle)
}

func (r *Router) readMultipartFile(req *http.Request) (store.CreateFileItemInput, error) {
	reader, err := req.MultipartReader()
	if err != nil {
		return store.CreateFileItemInput{}, err
	}

	var (
		sourceDeviceID   string
		sourceDeviceName string
		result           store.CreateFileItemInput
		foundFile        bool
	)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return store.CreateFileItemInput{}, err
		}

		switch part.FormName() {
		case "file":
			if foundFile {
				_ = part.Close()
				return store.CreateFileItemInput{}, fmt.Errorf("only one file part is supported")
			}

			result, err = r.saveMultipartFile(req.Context(), part)
			_ = part.Close()
			if err != nil {
				return store.CreateFileItemInput{}, err
			}
			foundFile = true
		case "source_device_id":
			rawValue, _ := io.ReadAll(part)
			sourceDeviceID = string(rawValue)
			_ = part.Close()
		case "source_device_name":
			rawValue, _ := io.ReadAll(part)
			sourceDeviceName = string(rawValue)
			_ = part.Close()
		default:
			_, _ = io.Copy(io.Discard, part)
			_ = part.Close()
		}
	}

	if !foundFile {
		return store.CreateFileItemInput{}, fmt.Errorf("multipart request must include a file field")
	}

	result.SourceDeviceID = strings.TrimSpace(sourceDeviceID)
	result.SourceDeviceName = strings.TrimSpace(sourceDeviceName)
	if r.cleaner != nil {
		settings, err := r.cleaner.GetSettings(req.Context())
		if err == nil {
			result.ExpiresAt = time.Now().UTC().Add(time.Duration(settings.TTLHours) * time.Hour).Format(time.RFC3339)
		}
	}

	return result, nil
}

func (r *Router) saveMultipartFile(ctx context.Context, part *multipart.Part) (store.CreateFileItemInput, error) {
	tempDir := filepath.Join(r.config.Storage.DataDir, "uploads", "tmp")
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		return store.CreateFileItemInput{}, err
	}

	tempFile, err := os.CreateTemp(tempDir, "clipbridge-upload-*")
	if err != nil {
		return store.CreateFileItemInput{}, err
	}
	defer tempFile.Close()

	hash := sha256.New()
	headerBytes := make([]byte, 0, 512)
	buffer := make([]byte, 32*1024)
	var sizeBytes int64

	for {
		readCount, readErr := part.Read(buffer)
		if readCount > 0 {
			chunk := buffer[:readCount]
			sizeBytes += int64(readCount)
			if len(headerBytes) < 512 {
				needed := 512 - len(headerBytes)
				if readCount < needed {
					needed = readCount
				}
				headerBytes = append(headerBytes, chunk[:needed]...)
			}
			if _, err := hash.Write(chunk); err != nil {
				return store.CreateFileItemInput{}, err
			}
			if _, err := tempFile.Write(chunk); err != nil {
				return store.CreateFileItemInput{}, err
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return store.CreateFileItemInput{}, readErr
		}
	}

	mimeType := http.DetectContentType(headerBytes)
	itemType := "file"
	if strings.HasPrefix(mimeType, "image/") {
		itemType = "image"
	}

	if err := r.validateBinary(ctx, itemType, sizeBytes); err != nil {
		_ = os.Remove(tempFile.Name())
		return store.CreateFileItemInput{}, err
	}

	sha256Hex := hex.EncodeToString(hash.Sum(nil))
	finalDir := filepath.Join(r.config.Storage.DataDir, "uploads", sha256Hex[:2])
	if err := os.MkdirAll(finalDir, 0o755); err != nil {
		return store.CreateFileItemInput{}, err
	}

	filename := sanitizeFilename(part.FileName())
	if filename == "" {
		filename = sha256Hex
	}
	finalPath := filepath.Join(finalDir, sha256Hex+"-"+filename)
	if _, err := os.Stat(finalPath); err == nil {
		_ = os.Remove(tempFile.Name())
	} else if err := os.Rename(tempFile.Name(), finalPath); err != nil {
		return store.CreateFileItemInput{}, err
	}

	return store.CreateFileItemInput{
		ItemType:  itemType,
		LocalPath: finalPath,
		Filename:  filename,
		MIMEType:  mimeType,
		SHA256:    sha256Hex,
		SizeBytes: sizeBytes,
	}, nil
}

func (r *Router) validateLink(ctx context.Context, rawURL string) error {
	limits := r.currentLimits(ctx)
	sizeBytes := len([]byte(rawURL))
	if sizeBytes < limits.MinLinkBytes {
		return fmt.Errorf("link must be at least %d bytes", limits.MinLinkBytes)
	}
	if sizeBytes > limits.MaxLinkBytes {
		return fmt.Errorf("link must be at most %d bytes", limits.MaxLinkBytes)
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("link must be a valid absolute URL")
	}

	return nil
}

func (r *Router) validateBinary(ctx context.Context, itemType string, sizeBytes int64) error {
	limits := r.currentLimits(ctx)

	switch itemType {
	case "image":
		if sizeBytes < int64(limits.MinImageBytes) {
			return fmt.Errorf("image must be at least %d bytes", limits.MinImageBytes)
		}
		if sizeBytes > int64(limits.MaxImageBytes) {
			return fmt.Errorf("image must be at most %d bytes", limits.MaxImageBytes)
		}
	case "file":
		if sizeBytes < int64(limits.MinFileBytes) {
			return fmt.Errorf("file must be at least %d bytes", limits.MinFileBytes)
		}
		if sizeBytes > int64(limits.MaxFileBytes) {
			return fmt.Errorf("file must be at most %d bytes", limits.MaxFileBytes)
		}
	}

	return nil
}

func mapMultipartError(err error) (int, string) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return http.StatusRequestEntityTooLarge, "request body is too large"
	}
	if strings.Contains(err.Error(), "multipart request must include a file field") {
		return http.StatusBadRequest, err.Error()
	}
	if strings.Contains(err.Error(), "only one file part is supported") {
		return http.StatusBadRequest, err.Error()
	}
	if errors.Is(err, io.EOF) {
		return http.StatusBadRequest, "multipart body must not be empty"
	}
	switch {
	case strings.Contains(err.Error(), "image must be"),
		strings.Contains(err.Error(), "file must be"):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusBadRequest, "multipart upload is invalid"
	}
}

func sanitizeFilename(filename string) string {
	filename = filepath.Base(strings.TrimSpace(filename))
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	return filename
}

func chooseContentType(mimeType string) string {
	if strings.TrimSpace(mimeType) == "" {
		return "application/octet-stream"
	}
	return mimeType
}
