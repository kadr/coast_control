package logger

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

const myTemplate = "{{datetime}} {{level}} [{{caller}}] {{message}}\n"

type MyLogger struct {
	slog.Logger
}

func New() *MyLogger {
	h := handler.NewConsoleHandler(slog.AllLevels)
	h.SetFormatter(slog.NewTextFormatter(myTemplate))
	h.TextFormatter().EnableColor = true
	log := slog.NewWithHandlers(h)
	return &MyLogger{Logger: *log}
}

func (l MyLogger) Info(message ...any) {
	l.Logger.Info(message)
}
func (l MyLogger) Infof(format string, args ...any) {
	l.Logger.Infof(format, args)
}
func (l MyLogger) Warn(message ...any) {
	l.Logger.Warn(message)
}
func (l MyLogger) Warnf(format string, args ...any) {
	l.Logger.Warnf(format, args)
}
func (l MyLogger) Error(message ...any) {
	l.Logger.Error(message)
}
func (l MyLogger) Errorf(format string, args ...any) {
	l.Logger.Errorf(format, args)
}
func (l MyLogger) Debug(message ...any) {
	l.Logger.Debug(message)
}
func (l MyLogger) Debugf(format string, args ...any) {
	l.Logger.Debugf(format, args)
}
func (l MyLogger) Fatal(message ...any) {
	l.Logger.Fatal(message)
}
func (l MyLogger) Fatalf(format string, args ...any) {
	l.Logger.Fatalf(format, args)
}
