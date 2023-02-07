package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type logging[C any, R any] struct {
	logger *logrus.Entry
}

func CreateLogging[C any, R any](logger *logrus.Entry) logging[C, R] {
	return logging[C, R]{
		logger: logger,
	}
}

func (d logging[C, R]) Handle(ctx context.Context, typeField string, command C, handler CommandHandler[C]) (err error) {
	logger := d.logger.WithFields(logrus.Fields{
		typeField:           generateActionName(command),
		typeField + "_body": fmt.Sprintf("%#v", command),
	})
	logger.Debug("Executing " + typeField)
	err = handler.Handle(ctx, command)
	if err == nil {
		logger.Info(typeField + " executed successfully")
	} else {
		logger.WithError(err).Error("Failed to execute " + typeField)
	}
	return err
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

func (d logging[C, R]) HandleQuery(ctx context.Context, command C, handler QueryHandler[C, R]) (result R, err error) {
	logger := d.logger.WithFields(logrus.Fields{
		"query":      generateActionName(command),
		"query_body": fmt.Sprintf("%#v", command),
	})

	logger.Debug("Executing query")
	r, err := handler.Handle(ctx, command)

	if err == nil {
		logger.Info("Query executed successfully")
	} else {
		logger.WithError(err).Error("Failed to execute query")
	}

	return r, err
}
