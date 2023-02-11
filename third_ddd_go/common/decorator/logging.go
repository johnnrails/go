package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type logging struct {
	logger *logrus.Entry
}

func CreateLogging(logger *logrus.Entry) logging {
	return logging{logger}
}

func (d logging) Handle(ctx context.Context, cmd interface{}, handler CommandHandler) (err error) {
	logger := d.logger.WithFields(logrus.Fields{
		"command":      generateActionName(cmd),
		"command_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Debug("Executing command")
	err = handler.Handle(ctx, cmd)
	if err == nil {
		logger.Info("command executed successfully")
	} else {
		logger.WithError(err).Error("Failed to execute command")
	}
	return err
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

func (d logging) HandleQuery(ctx context.Context, cmd interface{}, handler QueryHandler) (result interface{}, err error) {
	logger := d.logger.WithFields(logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	})

	logger.Debug("Executing query")
	r, err := handler.Handle(ctx, cmd)

	if err == nil {
		logger.Info("Query executed successfully")
	} else {
		logger.WithError(err).Error("Failed to execute query")
	}

	return r, err
}
