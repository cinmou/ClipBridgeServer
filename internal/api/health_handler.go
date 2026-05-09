// SPDX-License-Identifier: GPL-3.0-only

package api

import "net/http"

const version = "0.1.0"

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
