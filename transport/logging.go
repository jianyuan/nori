package transport

import "github.com/Sirupsen/logrus"

var log = logrus.New()

func init() {
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
}
