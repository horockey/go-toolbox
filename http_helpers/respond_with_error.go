package http_helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Respond with given error.
// Response body will have struct of errorMessage (see above).
func RespondWithErr(w http.ResponseWriter, statusCode int, err error) error {
	if err == nil {
		err = errors.New("")
	}
	b, err := json.Marshal(errorMessage{err.Error()})
	if err != nil {
		err = fmt.Errorf("wrapping errror message to json: %w", err)
		return err
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}
