package provider

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/genai"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		kind       ErrorKind
		statusCode int
	}{
		{name: "rate limit", err: &genai.APIError{Code: 429, Message: "rate limit"}, kind: ErrorKindRateLimit, statusCode: 429},
		{name: "authentication", err: &genai.APIError{Code: 401}, kind: ErrorKindAuthentication, statusCode: 401},
		{name: "invalid request", err: &genai.APIError{Code: 400}, kind: ErrorKindInvalidRequest, statusCode: 400},
		{name: "server", err: &genai.APIError{Code: 503}, kind: ErrorKindServer, statusCode: 503},
		{name: "timeout", err: context.DeadlineExceeded, kind: ErrorKindTimeout, statusCode: 504},
		{name: "unknown", err: errors.New("unexpected provider failure"), kind: ErrorKindUnknown, statusCode: 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifiedErr := ClassifyError("gemini", tt.err)
			if classifiedErr.Kind != tt.kind || classifiedErr.StatusCode != tt.statusCode {
				t.Fatalf("ClassifyError() = kind %q, status %d; want kind %q, status %d", classifiedErr.Kind, classifiedErr.StatusCode, tt.kind, tt.statusCode)
			}
		})
	}
}
