package users

import (
	"context"
	"fmt"

	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/storage"
	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/trace"
)

type Service struct {
	storage storage.UserStorer
}

func (s Service) Execute(ctx context.Context, req *Request) error {
	ctx, span := trace.NewSpan(ctx, "userService.execute", nil)
	defer span.End()

	if err := s.storage.Insert(ctx, storage.User{Name: req.Name}); err != nil {
		return fmt.Errorf("create: unable to store: %w", err)
	}

	return nil
}