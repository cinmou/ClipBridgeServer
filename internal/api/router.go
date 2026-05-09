// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/auth"
	"github.com/cinmou/ClipBridgeServer/internal/cleanup"
	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/store"
)

const pairingCodeTTL = 5 * time.Minute

// Router wires HTTP handlers to the concrete dependencies assembled in main.
// Keeping the dependency graph explicit makes the service easier to reason
// about as a portable single-binary deployment target.
type Router struct {
	store   *store.SQLiteStore
	config  *config.Config
	cleaner *cleanup.Service
	webUI   http.Handler
}

// NewRouter builds the route table for the current server phase.
func NewRouter(dbStore *store.SQLiteStore, cfg *config.Config, cleanerService *cleanup.Service, webUI http.Handler) http.Handler {
	router := &Router{
		store:   dbStore,
		config:  cfg,
		cleaner: cleanerService,
		webUI:   webUI,
	}

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/api/health", healthHandler)
	publicMux.HandleFunc("/api/auth/pair", router.handlePairDevice)

	clipboardMux := http.NewServeMux()
	clipboardMux.HandleFunc("/api/clipboard/text", router.handleClipboardText)
	clipboardMux.HandleFunc("/api/clipboard/latest", router.handleClipboardLatest)
	clipboardMux.HandleFunc("/api/clipboard/history", router.handleClipboardHistory)
	clipboardMux.HandleFunc("/api/clipboard/items/", router.handleClipboardItemRoutes)
	clipboardMux.HandleFunc("/api/favorites", router.handleFavorites)
	clipboardMux.HandleFunc("/api/categories", router.handleCategories)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/api/auth/pairing-codes", router.handlePairingCodes)
	adminMux.HandleFunc("/api/auth/devices", router.handleDevices)
	adminMux.HandleFunc("/api/auth/devices/", router.handleDeviceByID)
	adminMux.HandleFunc("/api/admin/cleanup/run", router.handleAdminCleanupRun)
	adminMux.HandleFunc("/api/admin/cleanup/status", router.handleAdminCleanupStatus)
	adminMux.HandleFunc("/api/admin/storage/status", router.handleAdminStorageStatus)
	adminMux.HandleFunc("/api/settings/cleanup", router.handleCleanupSettings)

	return router.loggingMiddleware(
		router.corsMiddleware(
			router.requestSizeMiddleware(
				router.routeAPIMiddleware(
					publicMux,
					router.adminAuthMiddleware(adminMux),
					router.clipboardAuthMiddleware(clipboardMux),
				),
			),
		),
	)
}

func (r *Router) routeAPIMiddleware(publicMux http.Handler, adminMux http.Handler, clipboardMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.URL.Path == "/api/health":
			publicMux.ServeHTTP(w, req)
		case req.URL.Path == "/api/auth/pair":
			publicMux.ServeHTTP(w, req)
		case strings.HasPrefix(req.URL.Path, "/api/auth/"):
			adminMux.ServeHTTP(w, req)
		case strings.HasPrefix(req.URL.Path, "/api/admin/"):
			adminMux.ServeHTTP(w, req)
		case strings.HasPrefix(req.URL.Path, "/api/settings/"):
			adminMux.ServeHTTP(w, req)
		case strings.HasPrefix(req.URL.Path, "/api/clipboard/"):
			clipboardMux.ServeHTTP(w, req)
		case req.URL.Path == "/api/favorites":
			clipboardMux.ServeHTTP(w, req)
		case req.URL.Path == "/api/categories":
			clipboardMux.ServeHTTP(w, req)
		default:
			r.webUI.ServeHTTP(w, req)
		}
	})
}

func (r *Router) adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := auth.ExtractBearerToken(req)
		if !ok || token != r.config.Auth.Token {
			writeJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Router) clipboardAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := auth.ExtractBearerToken(req)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if token == r.config.Auth.Token {
			next.ServeHTTP(w, req)
			return
		}

		if _, err := r.store.AuthenticateDeviceToken(req.Context(), token); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "authenticate device token failed")
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Router) requestSizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// One global request-size limit keeps the API on a predictable safety
		// baseline before any handler starts decoding JSON.
		req.Body = http.MaxBytesReader(w, req.Body, int64(r.config.Limits.MaxRequestBytes))
		next.ServeHTTP(w, req)
	})
}

func (r *Router) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if allowedOrigin := allowOrigin(req.Header.Get("Origin")); allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		}

		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, req)

		log.Printf("%s %s %d %s", req.Method, req.URL.Path, recorder.statusCode, time.Since(startTime))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func allowOrigin(origin string) string {
	switch origin {
	case "http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:4173",
		"http://127.0.0.1:4173",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:8787",
		"http://127.0.0.1:8787":
		return origin
	default:
		return ""
	}
}
