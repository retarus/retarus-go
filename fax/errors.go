package fax

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("bad request: Client authorization is missing")
	ErrAuthFailure         = errors.New("authentication failure: bad or missing authentication")
	ErrNotFound            = errors.New("not found: No job/recipient report available for the given jobId")
	ErrConflict            = errors.New("conflict: Duplicate job")
	ErrInternalServerError = errors.New("internal server error: Cannot accept job, cannot query jobReport, cannot list jobs, cannot query recipient report, cannot apply transliteration in the send job")
	ErrServiceUnavailable  = errors.New("service unavailable: The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.")
	ErrUnknown             = errors.New("unknown error: Server signals that there was an unknown problem, most likely with the backend adaptor")
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
