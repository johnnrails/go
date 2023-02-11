package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type commandMetrics struct {
	client MetricsClient
}

func (d commandMetrics) Handle(ctx context.Context, cmd interface{}, handler CommandHandler) (err error) {
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

type queryMetrics struct {
	client MetricsClient
}

func (d queryMetrics) Handle(ctx context.Context, query interface{}, handler QueryHandler) (result interface{}, err error) {
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
