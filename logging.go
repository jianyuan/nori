package nori

import (
	"github.com/Sirupsen/logrus"
	"github.com/jianyuan/nori/log"
	"golang.org/x/net/context"
)

func configureLogger(ctx context.Context, config *Configuration) (context.Context, error) {
	if logger, ok := log.FromContext(ctx).(*logrus.Logger); ok {
		logger.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
		ctx = log.NewContext(ctx, logger)
	}
	return ctx, nil
}
