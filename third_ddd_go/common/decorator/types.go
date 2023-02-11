package decorator

import "context"

type CommandHandler interface {
	Handle(ctx context.Context, cmd any) error
}

type QueryHandler interface {
	Handle(ctx context.Context, q any) (any, error)
}

type MetricsClient interface {
	Inc(key string, value int)
}
