package logs

import "github.com/sirupsen/logrus"

func Init() {
	SetFormatter(logrus.StandardLogger())

	logrus.SetLevel(logrus.DebugLevel)
}

// This funcion is used around the application to config
// the formatter for the logger
func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
}
