package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type StructuredLogger struct {
	logger *logrus.Logger
}

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	logFields := logrus.Fields{}
	logFields["http_method"] = r.Method
	logFields["remote_addr"] = r.RemoteAddr
	logFields["uri"] = r.RequestURI
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	entry := &StructuredLoggerEntry{FieldLogger: logrus.NewEntry(l.logger)}
	entry.FieldLogger = entry.FieldLogger.WithFields(logFields)
	entry.FieldLogger.Info("Request started")

	return entry
}

type StructuredLoggerEntry struct {
	FieldLogger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.FieldLogger = l.FieldLogger.WithFields(logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed":      elapsed.Round(time.Millisecond / 100).String(),
	})
	l.FieldLogger.Info("Request completed")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.FieldLogger = l.FieldLogger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	return middleware.GetLogEntry(r).(*StructuredLoggerEntry).FieldLogger
}
