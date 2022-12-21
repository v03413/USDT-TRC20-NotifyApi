package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var log *logrus.Logger

const defaultOutput = "/var/log/main.log"

func Init(logLevel, output string) {
	log = logrus.New()
	var level = logrus.DebugLevel

	switch {
	case logLevel == "debug":
		level = logrus.DebugLevel
	case logLevel == "info":
		level = logrus.InfoLevel
	case logLevel == "error":
		level = logrus.ErrorLevel
	default:
		level = logrus.DebugLevel
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	logFile := defaultOutput
	if output != "" {
		logFile = output
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {

		log.Error(err.Error())
	}

	if level == logrus.DebugLevel {
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		log.SetOutput(file)
	}
}

func Println(v ...interface{}) {

	log.Infoln(v)
}

func Error(v ...interface{}) {

	log.Errorln(v)
}

func Debug(v ...interface{}) {

	log.Debugln(v)
}

func Fatal(v ...interface{}) {

	log.Fatalln(v)
}
