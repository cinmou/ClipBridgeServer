// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type dataResponse struct {
	Data any `json:"data"`
}

type apiError struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSONData(w http.ResponseWriter, statusCode int, payload any) {
	writeJSON(w, statusCode, dataResponse{Data: payload})
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, errorResponse{Error: apiError{Message: message}})
}

func writeMethodNotAllowed(w http.ResponseWriter, methodsAllowed ...string) {
	w.Header().Set("Allow", strings.Join(methodsAllowed, ", "))
	writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
}
