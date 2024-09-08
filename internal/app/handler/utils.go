package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpHandler = func(w http.ResponseWriter, r *http.Request)

// decode decodes a JSON request body into a provided struct
func decodeJsonReq[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// SendJson sends a JSON response to the client
func sendJson(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
