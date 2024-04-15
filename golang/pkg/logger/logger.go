package logger

import "go.uber.org/zap"

func NewLogger() Logger {
	zapLogger, _ := zap.NewProduction()
	return &logger{
		lg: zapLogger,
	}
}

type Logger interface {
	Info(message string)
	Error(message string, err error)
	Fatal(message string, err error)
}

type logger struct {
	lg *zap.Logger
}

func (l *logger) Info(message string) {
	l.lg.Info(message)
}

func (l *logger) Error(message string, err error) {
	l.lg.Error(message, zap.Error(err))
}

func (l *logger) Fatal(message string, err error) {
	l.lg.Fatal(message, zap.Error(err))
}
