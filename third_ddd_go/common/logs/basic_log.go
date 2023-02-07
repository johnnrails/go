package logs

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func Basic_log() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	log.Println("This is my log message.")
}

func configLogrus() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))

	if err != nil {
		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)
}

func Basic_Logrus() {

	configLogrus()

	e := echo.New()
	e.Use(loggingMiddleware)

	e.GET("/test/logrus", func(ctx echo.Context) error {
		a := ctx.QueryParam("a")

		logEntry := logrus.WithField("a", a)
		logEntry.Debug("parsing param \"a\"")

		if _, err := strconv.Atoi(a); err != nil {
			logEntry.Debug("unable to parse \"a\" param")
			return ctx.String(http.StatusBadRequest, "not ol")
		}

		logEntry.Debug("parsed \"a\" param")

		return ctx.String(http.StatusOK, "ok")
	})

	e.Start(":8080")
}

func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()

		res := next(ctx)

		logrus.WithFields(logrus.Fields{
			"method":     ctx.Request().Method,
			"path":       ctx.Path(),
			"status":     ctx.Response().Status,
			"latency_ns": time.Since(start).Nanoseconds(),
		}).Info("request")

		return res
	}
}
