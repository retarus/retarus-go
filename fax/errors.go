package fax

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("client Error: Authorization header is missing or incorrect")
	ErrAuthFailure         = errors.New("authentication Failed: Invalid or missing credentials")
	ErrNotFound            = errors.New("resource Not Found: No report exists for the specified job ID")
	ErrConflict            = errors.New("conflict: A job with the same identifier already exists")
	ErrInternalServerError = errors.New("internal Server Error: Various possible issues including job acceptance, report queries, job listings, and transliteration")
	ErrServiceUnavailable  = errors.New("service Unavailable: Server is temporarily overloaded or under maintenance")
	ErrUnknown             = errors.New("unknown Error: An unspecified issue occurred, possibly related to the backend adaptor")
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

	msgStr := string(msg)

	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrBadRequest, msgStr)
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrAuthFailure, msgStr)
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, msgStr)
	case http.StatusConflict:
		return fmt.Errorf("%w: %s", ErrConflict, msgStr)
	case http.StatusInternalServerError:
		return fmt.Errorf("%w: %s", ErrInternalServerError, msgStr)
	case http.StatusServiceUnavailable:
		return fmt.Errorf("%w: %s", ErrServiceUnavailable, msgStr)
	}

	return fmt.Errorf("%w: %s", ErrUnknown, msgStr)
}
