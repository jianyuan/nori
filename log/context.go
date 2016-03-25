package log

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type loggerKey struct{}

func NewContext(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) logrus.FieldLogger {
	if logger, ok := ctx.Value(loggerKey{}).(*logrus.Logger); ok {
		return logger
	}
	return logrus.StandardLogger()
}
