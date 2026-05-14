// SPDX-License-Identifier: GPL-3.0-only

package api

import "net/http"

const version = "v0.2.0-beta.1"

type healthResponse struct {
	OK      bool   `json:"ok"`
	Version string `json:"version"`
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSONData(w, http.StatusOK, healthResponse{
		OK:      true,
		Version: version,
	})
}
