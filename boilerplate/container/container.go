package container

import (
	"context"

	"github.com/vardius/gocontainer"
)

const containerKey int = 1

func ContextWithContainer(ctx context.Context, c gocontainer.Container) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, containerKey, c)
}

func FromContext(ctx context.Context) (gocontainer.Container, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(containerKey).(gocontainer.Container)
	return c, ok
}
