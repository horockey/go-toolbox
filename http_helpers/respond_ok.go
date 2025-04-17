package http_helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

// Respond with status 200.
// If marshaling error occures, response with error and code 500 will be sent.
func RespondOK(w http.ResponseWriter, data any) error {
	if data == nil {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}
