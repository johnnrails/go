package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/trace"
)

type Request struct {
	Name string `json:"name"`
}

func (r *Request) validate(ctx context.Context, body io.Reader) error {
	_, span := trace.NewSpan(ctx, "userRequest.validate", nil)
	defer span.End()
	
	if err := json.NewDecoder(body).Decode(r); err != nil {
		return fmt.Errorf("validate: malformed body")
	}

	if r.Name == "" {
		return fmt.Errorf("validate: invalid request")
	}

	return nil
}