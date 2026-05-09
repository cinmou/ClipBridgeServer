// SPDX-License-Identifier: GPL-3.0-only

package webui

import (
	"embed"
	"io/fs"
	"net/http"
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
		if req.URL.Path == "/" {
			indexHTML, readErr := fs.ReadFile(staticFS, "index.html")
			if readErr != nil {
				http.Error(w, "embedded index.html is missing", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(indexHTML)
			return
		}

		req = req.Clone(req.Context())
		fileServer.ServeHTTP(w, req)
	})
}
