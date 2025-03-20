package log

import (
	"go.uber.org/zap"
	stdlog "log"
)

var _ Logger = (*noopLogger)(nil)

var (
	Log Logger
	Nop Logger = &noopLogger{}

	ZapLog *zap.Logger
)

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)

	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)
}

func init() {
	var (
		err       error
		zapLogger *zap.Logger
	)
	zapLogger, err = zap.NewDevelopment()
	if err != nil {
		stdlog.Fatalf("can't initialize zap logger: %v", err)
	}

	ZapLog = zapLogger

	slogger := zapLogger.Sugar()
	Log = slogger
}

type noopLogger struct {
}

func (n *noopLogger) Debug(args ...any) {
}

func (n *noopLogger) Info(args ...any) {
}

func (n *noopLogger) Warn(args ...any) {
}

func (n *noopLogger) Error(args ...any) {
}

func (n *noopLogger) Fatal(args ...any) {
}

func (n *noopLogger) Debugf(format string, args ...any) {
}

func (n *noopLogger) Infof(format string, args ...any) {
}

func (n *noopLogger) Warnf(format string, args ...any) {
}

func (n *noopLogger) Errorf(format string, args ...any) {
}

func (n *noopLogger) Fatalf(format string, args ...any) {
}

func (n *noopLogger) Debugw(msg string, keysAndValues ...any) {
}

func (n *noopLogger) Infow(msg string, keysAndValues ...any) {
}

func (n *noopLogger) Warnw(msg string, keysAndValues ...any) {
}

func (n *noopLogger) Errorw(msg string, keysAndValues ...any) {
}

func (n *noopLogger) Fatalw(msg string, keysAndValues ...any) {
}
