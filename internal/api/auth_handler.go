// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/store"
)

type pairingCodeResponse struct {
	PairingCode string `json:"pairing_code"`
	ExpiresAt   string `json:"expires_at"`
	PairingURI  string `json:"pairing_uri"`
}

type pairDeviceRequest struct {
	PairingCode string `json:"pairing_code"`
	DeviceName  string `json:"device_name"`
}

type pairDeviceResponse struct {
	DeviceToken string         `json:"device_token"`
	Device      deviceResponse `json:"device"`
}

type deviceResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	LastSeenAt string `json:"last_seen_at,omitempty"`
	RevokedAt  string `json:"revoked_at,omitempty"`
}

type devicesResponse struct {
	Items []deviceResponse `json:"items"`
}

func (r *Router) handlePairingCodes(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	expiresAt := time.Now().UTC().Add(pairingCodeTTL)
	pairingCode, err := r.store.CreatePairingCode(req.Context(), expiresAt)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "create pairing code failed")
		return
	}

	writeJSONData(w, http.StatusCreated, pairingCodeResponse{
		PairingCode: pairingCode.Code,
		ExpiresAt:   pairingCode.ExpiresAt,
		PairingURI:  buildPairingURI(req, pairingCode.Code),
	})
}

func (r *Router) handlePairDevice(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	var payload pairDeviceRequest
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

	if strings.TrimSpace(payload.PairingCode) == "" {
		writeJSONError(w, http.StatusBadRequest, "pairing_code must not be empty")
		return
	}

	device, deviceToken, err := r.store.ExchangePairingCode(req.Context(), payload.PairingCode, payload.DeviceName)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusBadRequest, "pairing code is invalid")
		case errors.Is(err, store.ErrAlreadyUsed):
			writeJSONError(w, http.StatusBadRequest, "pairing code has already been used")
		case errors.Is(err, store.ErrExpired):
			writeJSONError(w, http.StatusBadRequest, "pairing code has expired")
		default:
			writeJSONError(w, http.StatusInternalServerError, "pair device failed")
		}
		return
	}

	writeJSONData(w, http.StatusCreated, pairDeviceResponse{
		DeviceToken: deviceToken,
		Device:      toDeviceResponse(device),
	})
}

func (r *Router) handleDevices(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	devices, err := r.store.ListDevices(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "load devices failed")
		return
	}

	response := devicesResponse{
		Items: make([]deviceResponse, 0, len(devices)),
	}
	for _, device := range devices {
		response.Items = append(response.Items, toDeviceResponse(&device))
	}

	writeJSONData(w, http.StatusOK, response)
}

func (r *Router) handleDeviceByID(w http.ResponseWriter, req *http.Request) {
	id, err := parseDeviceID(req.URL.Path)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	if req.Method != http.MethodDelete {
		writeMethodNotAllowed(w, http.MethodDelete)
		return
	}

	if err := r.store.RevokeDevice(req.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "device not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "revoke device failed")
		return
	}

	writeJSONData(w, http.StatusOK, map[string]any{
		"revoked": true,
		"id":      id,
	})
}

func toDeviceResponse(device *store.Device) deviceResponse {
	return deviceResponse{
		ID:         device.ID,
		Name:       device.Name,
		CreatedAt:  device.CreatedAt,
		LastSeenAt: device.LastSeenAt,
		RevokedAt:  device.RevokedAt,
	}
}

func parseDeviceID(path string) (int64, error) {
	const prefix = "/api/auth/devices/"

	if !strings.HasPrefix(path, prefix) {
		return 0, fmt.Errorf("device not found")
	}

	idText := strings.TrimPrefix(path, prefix)
	if idText == "" || strings.Contains(idText, "/") {
		return 0, fmt.Errorf("device not found")
	}

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("device not found")
	}

	return id, nil
}

func buildPairingURI(req *http.Request, pairingCode string) string {
	serverURL := "http://" + req.Host
	return "clipbridge://pair?server_url=" + serverURL + "&pairing_code=" + pairingCode
}
