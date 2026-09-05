package provider

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genai"
)

type ErrorKind string

const (
	ErrorKindAuthentication ErrorKind = "authentication"
	ErrorKindAuthorization  ErrorKind = "authorization"
	ErrorKindInvalidRequest ErrorKind = "invalid_request"
	ErrorKindRateLimit      ErrorKind = "rate_limit"
	ErrorKindNotFound       ErrorKind = "not_found"
	ErrorKindTimeout        ErrorKind = "timeout"
	ErrorKindServer         ErrorKind = "server"
	ErrorKindCanceled       ErrorKind = "canceled"
	ErrorKindUnknown        ErrorKind = "unknown"
)

type Error struct {
	Provider   string
	Kind       ErrorKind
	StatusCode int
	Err        error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s provider %s error: %v", e.Provider, e.Kind, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func ClassifyError(providerName string, err error) *Error {
	if err == nil {
		return nil
	}

	providerErr := &Error{
		Provider:   providerName,
		Kind:       ErrorKindUnknown,
		StatusCode: 500,
		Err:        err,
	}

	switch {
	case errors.Is(err, context.Canceled):
		providerErr.Kind = ErrorKindCanceled
		providerErr.StatusCode = 499
	case errors.Is(err, context.DeadlineExceeded):
		providerErr.Kind = ErrorKindTimeout
		providerErr.StatusCode = 504
	default:
		if statusCode, ok := apiErrorStatusCode(err); ok {
			providerErr.StatusCode = statusCode
			providerErr.Kind = classifyStatusCode(statusCode)
		}
	}

	return providerErr
}

func apiErrorStatusCode(err error) (int, bool) {
	var apiErr *genai.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code, apiErr.Code != 0
	}

	var apiErrValue genai.APIError
	if errors.As(err, &apiErrValue) {
		return apiErrValue.Code, apiErrValue.Code != 0
	}

	return 0, false
}

func classifyStatusCode(statusCode int) ErrorKind {
	switch statusCode {
	case 400, 422:
		return ErrorKindInvalidRequest
	case 401:
		return ErrorKindAuthentication
	case 403:
		return ErrorKindAuthorization
	case 404:
		return ErrorKindNotFound
	case 429:
		return ErrorKindRateLimit
	case 500, 502, 503, 504:
		return ErrorKindServer
	default:
		if statusCode >= 400 && statusCode < 500 {
			return ErrorKindInvalidRequest
		}
		if statusCode >= 500 {
			return ErrorKindServer
		}
		return ErrorKindUnknown
	}
}
