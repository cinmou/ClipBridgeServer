// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cinmou/ClipBridgeServer/internal/store"
)

type appSettingsPatchRequest struct {
	AdminToken string `json:"admin_token"`
}

func (r *Router) handleSettings(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		settings, err := r.loadAppSettings(req.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "load settings failed")
			return
		}
		writeJSONData(w, http.StatusOK, settings)
	case http.MethodPatch:
		var payload appSettingsPatchRequest
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

		if strings.TrimSpace(payload.AdminToken) == "" {
			writeJSONError(w, http.StatusBadRequest, "admin_token must not be empty")
			return
		}
		if err := r.store.SaveAdminToken(req.Context(), strings.TrimSpace(payload.AdminToken)); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "save settings failed")
			return
		}

		settings, err := r.loadAppSettings(req.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "load settings failed")
			return
		}
		writeJSONData(w, http.StatusOK, settings)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPatch)
	}
}

func (r *Router) handleSettingsLimits(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		writeJSONData(w, http.StatusOK, r.currentLimits(req.Context()))
	case http.MethodPatch:
		var payload store.LimitsSettings
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
		if err := payload.Validate(); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := r.store.SaveLimitsSettings(req.Context(), payload); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "save limits failed")
			return
		}
		writeJSONData(w, http.StatusOK, payload)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPatch)
	}
}

func (r *Router) loadAppSettings(ctx context.Context) (store.AppSettings, error) {
	limits, err := r.store.LoadLimitsSettings(ctx, r.defaultLimits())
	if err != nil {
		return store.AppSettings{}, err
	}
	cleanupSettings, err := r.cleaner.GetSettings(ctx)
	if err != nil {
		return store.AppSettings{}, err
	}
	adminToken, err := r.store.LoadAdminToken(ctx, r.config.Auth.Token)
	if err != nil {
		return store.AppSettings{}, err
	}

	return store.AppSettings{
		AdminToken: adminToken,
		Limits:     limits,
		Cleanup:    cleanupSettings,
		WebDAV:     webdavSettingsOrZero(r, ctx),
		Startup: store.StartupSettings{
			Host:         r.config.Server.Host,
			Port:         r.config.Server.Port,
			DataDir:      r.config.Storage.DataDir,
			DatabasePath: r.config.Storage.DatabasePath,
		},
		RestartRequiredFields: []string{"host", "port", "data_dir", "database_path"},
	}, nil
}

func webdavSettingsOrZero(r *Router, ctx context.Context) store.WebDAVSettings {
	if r.webdav == nil {
		return store.WebDAVSettings{BasePath: "ClipBridgeServer"}
	}
	settings, err := r.webdav.GetSettings(ctx)
	if err != nil {
		return store.WebDAVSettings{BasePath: "ClipBridgeServer"}
	}
	return settings
}
