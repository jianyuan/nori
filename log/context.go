package log

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type loggerKey struct{}

func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return logger
	}
	return logrus.StandardLogger()
}
