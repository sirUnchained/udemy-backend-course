package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes a JSON response with the given status code and data
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "appliction/json")
	w.WriteHeader(status)
	// Encode data as JSON and write to response writer
	return json.NewEncoder(w).Encode(data)
}

// readJSON reads and decodes JSON from the request body into the provided data structure
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb or maximum allowed request body size
	// limit the size of the request body to prevent large payloads
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// create JSON decoder for the request body
	decoder := json.NewDecoder(r.Body)
	// reject requests with unknown fields in JSON
	decoder.DisallowUnknownFields()

	// Decode JSON into the provided data structure
	return decoder.Decode(data)
}

// writeJSONError writes a JSON error response with the given status code and message
func writeJSONError(w http.ResponseWriter, status int, msg string) error {
	// define envelope structure for error responses
	type envelope struct {
		Error string `json:"error"`
	}

	// use writeJSON to send the error response in consistent envelope format
	return writeJSON(w, status, &envelope{Error: msg})
}
