package http

import (
	"context"
	"net"
	"net/http"
)

type Adapter struct {
	httpServer *http.Server
}

func (a *Adapter) Start(ctx context.Context) error {
	a.httpServer.BaseContext = func(_ net.Listener) context.Context { return ctx }
	if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *Adapter) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
