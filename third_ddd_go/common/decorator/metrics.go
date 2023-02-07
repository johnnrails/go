package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type commandMetrics[C any] struct {
	client MetricsClient
}

func (d commandMetrics[C]) Handle(ctx context.Context, cmd C, handler CommandHandler[C]) (err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))

	defer func() {
		secondsSinceStart := int(time.Since(start).Seconds())
		d.client.Inc(fmt.Sprintf("commands.%s.duration", actionName), secondsSinceStart)
		if err == nil {
			d.client.Inc(fmt.Sprintf("commands.%s.success", actionName), 1)
		} else {
			d.client.Inc(fmt.Sprintf("commands.%s.failure", actionName), 1)
		}
	}()

	return handler.Handle(ctx, cmd)
}

type queryMetrics[C any, R any] struct {
	client MetricsClient
}

func (d queryMetrics[C, R]) Handle(ctx context.Context, query C, handler QueryHandler[C, R]) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(query))

	defer func() {
		secondsSinceStart := int(time.Since(start).Seconds())
		d.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), secondsSinceStart)
		if err == nil {
			d.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
		} else {
			d.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
		}
	}()

	return handler.Handle(ctx, query)
}
