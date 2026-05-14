// SPDX-License-Identifier: GPL-3.0-only

package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// distFS holds the built frontend bundle that ships inside the server binary.
// Embedding the files keeps deployment to one executable plus config and data.
//
//go:embed dist/*
var distFS embed.FS

// Handler serves the embedded Web UI bundle. The fallback to index.html allows
// the frontend to keep using simple browser navigation without a separate web
// server process.
func Handler() http.Handler {
	staticFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if shouldServeIndex(staticFS, req.URL.Path) {
			if err := serveIndexHTML(w, staticFS); err != nil {
				http.Error(w, "embedded index.html is missing", http.StatusInternalServerError)
			}
			return
		}

		req = req.Clone(req.Context())
		fileServer.ServeHTTP(w, req)
	})
}

func shouldServeIndex(staticFS fs.FS, requestPath string) bool {
	if requestPath == "/" || requestPath == "" {
		return true
	}

	cleanPath := strings.TrimPrefix(path.Clean(requestPath), "/")
	if cleanPath == "." || cleanPath == "" {
		return true
	}

	if strings.Contains(path.Base(cleanPath), ".") {
		return false
	}

	_, err := fs.Stat(staticFS, cleanPath)
	return err != nil
}

func serveIndexHTML(w http.ResponseWriter, staticFS fs.FS) error {
	indexHTML, err := fs.ReadFile(staticFS, "index.html")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
	return nil
}
