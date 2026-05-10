// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"encoding/json"
	"net/http"

	"github.com/cinmou/ClipBridgeServer/internal/store"
)

func (r *Router) handleWebDAVSettings(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		settings, err := r.webdav.GetSettings(req.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "load webdav settings failed")
			return
		}
		writeJSONData(w, http.StatusOK, settings)
	case http.MethodPatch:
		var payload store.WebDAVSettings
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
		settings, err := r.webdav.UpdateSettings(req.Context(), payload)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSONData(w, http.StatusOK, settings)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPatch)
	}
}

func (r *Router) handleAdminWebDAVTest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}
	status, err := r.webdav.TestConnection(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSONData(w, http.StatusOK, status)
}

func (r *Router) handleAdminWebDAVSync(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}
	status, err := r.webdav.RunSync(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSONData(w, http.StatusOK, status)
}

func (r *Router) handleAdminWebDAVStatus(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	status, err := r.webdav.GetStatus(req.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "load webdav status failed")
		return
	}
	writeJSONData(w, http.StatusOK, status)
}
