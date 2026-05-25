// Package httputil provides shared HTTP response helpers used by both the
// handler and middleware packages. It exists as a separate package to avoid
// an import cycle: middleware imports handler utilities, but handler already
// imports middleware — placing these helpers inside either package would create
// a circular dependency.
package httputil

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

// WriteError writes a JSON {"error": msg} response with the given status code.
func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

// WriteJSON writes v as JSON with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
