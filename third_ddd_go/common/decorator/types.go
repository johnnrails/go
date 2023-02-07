package decorator

import "context"

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

type MetricsClient interface {
	Inc(key string, value int)
}
