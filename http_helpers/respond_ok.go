package http_helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

// Respond with status 200.
// If marshaling error occures, response with error and code 500 will be sent.
func RespondOK(w http.ResponseWriter, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("marshaling json: %w", err)
		respErr := RespondWithErr(
			w,
			http.StatusInternalServerError,
			err,
		)
		if respErr != nil {
			err = errors.Join(err, respErr)
		}

		return fmt.Errorf("responding: %w", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}
