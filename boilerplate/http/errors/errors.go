package errors

import (
	"context"
	"net/http"

	"github.com/johnnrails/ddd_go/boilerplate/metadata"
)

type HttpError struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func NewHttpError(ctx context.Context, code int, err error) *HttpError {
	httpError := &HttpError{
		Code:    code,
		Message: http.StatusText(code),
	}

	if m := metadata.FromContext(ctx); m != nil {
		httpError.RequestID = m.TraceID
		m.Err = err
	}

	return httpError
}
