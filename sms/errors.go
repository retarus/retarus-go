package sms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func statusToError(statusCode int, body io.Reader) error {
	switch statusCode {
	case http.StatusCreated, http.StatusOK:
		return nil
	}

	msg, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	type errFormat struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	var errF errFormat
	if err := json.Unmarshal(msg, &errF); err == nil {
		return fmt.Errorf("%s: %s", http.StatusText(statusCode), errF.Message)
	}

	return fmt.Errorf("%s: %s", http.StatusText(statusCode), string(msg))
}
