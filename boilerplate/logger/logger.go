package logger

import (
	"context"
	"encoding/json"
	"log"

	"github.com/johnnrails/ddd_go/boilerplate/metadata"
	"github.com/vardius/golog"
)

var Logger golog.Logger

func SetFlags(flag int) {
	Logger.SetFlags(flag)
}

func SetVerbosity(verbosity golog.Verbose) {
	Logger.SetVerbosity(verbosity)
}

func Debug(ctx context.Context, v string) {
	Logger.Debug(ctx, messageWithMeta(ctx, v))
}

func Info(ctx context.Context, v string) {
	Logger.Info(ctx, messageWithMeta(ctx, v))
}

func Warning(ctx context.Context, v string) {
	Logger.Warning(ctx, messageWithMeta(ctx, v))
}

func Error(ctx context.Context, v string) {
	Logger.Error(ctx, messageWithMeta(ctx, v))
}

func Critical(ctx context.Context, v string) {
	Logger.Critical(ctx, messageWithMeta(ctx, v))
}

func Fatal(ctx context.Context, v string) {
	Logger.Fatal(ctx, messageWithMeta(ctx, v))
}

type messageStruct struct {
	Message string             `json:"message"`
	Meta    *metadata.Metadata `json:"meta"`
}

func messageWithMeta(ctx context.Context, v string) string {
	mtd := metadata.FromContext(ctx)
	s, _ := json.Marshal(messageStruct{
		Message: v,
		Meta:    mtd,
	})
	return string(s)
}

func init() {
	l := golog.New()
	l.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)

	Logger = l
}
