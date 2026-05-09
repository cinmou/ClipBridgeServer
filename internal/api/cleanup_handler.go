// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"encoding/json"
	"net/http"

	"github.com/cinmou/ClipBridgeServer/internal/store"
)

func (r *Router) handleAdminCleanupRun(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	status, err := r.cleaner.RunNow(req.Context(), "manual")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "run cleanup failed")
		return
	}

	writeJSONData(w, http.StatusOK, status)
}

func (r *Router) handleAdminCleanupStatus(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	status, err := r.cleaner.GetStatus(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "load cleanup status failed")
		return
	}

	writeJSONData(w, http.StatusOK, status)
}

func (r *Router) handleAdminStorageStatus(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	status, err := r.cleaner.GetStorageStatus(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "load storage status failed")
		return
	}

	writeJSONData(w, http.StatusOK, status)
}

func (r *Router) handleCleanupSettings(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		settings, err := r.cleaner.GetSettings(req.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "load cleanup settings failed")
			return
		}
		writeJSONData(w, http.StatusOK, settings)
	case http.MethodPatch:
		var payload store.CleanupSettings
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

		settings, err := r.cleaner.UpdateSettings(req.Context(), payload)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSONData(w, http.StatusOK, settings)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPatch)
	}
}
