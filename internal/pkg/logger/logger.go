package logger

import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

var (
	l *logrus.Logger
)

func InitLogger() {
	l = logrus.New()
	l.SetFormatter(&ecslogrus.Formatter{})
	l.SetLevel(logrus.DebugLevel)
}

func Debug(args ...interface{}) {
	l.Debugln(args...)
}

func Warn(args ...interface{}) {
	l.Warnln(args...)
}

func Info(args ...interface{}) {
	l.Infoln(args...)
}

func Error(args ...interface{}) {
	l.Errorln(args...)
}

func Fatal(args ...interface{}) {
	l.Fatalln(args...)
}

func Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}
