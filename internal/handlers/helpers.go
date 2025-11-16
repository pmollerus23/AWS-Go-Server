package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Validator is an interface for validating request payloads.
type Validator interface {
	Valid(ctx context.Context) map[string]string
}

// encode encodes a value as JSON and writes it to the response.
func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// decode decodes a request body into the provided type.
func decode[T any](r *http.Request, v *T) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// decodeValid decodes and validates a request body.
func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, err
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, errors.New("validation failed")
	}
	return v, nil, nil
}
