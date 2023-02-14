package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

var loggerKey string = "logger"

func InitLogger(ctx context.Context, level string) context.Context {
	if ctx.Value(loggerKey) != nil {
		return ctx
	}
	return context.WithValue(ctx, loggerKey, logrus.WithFields(
		logrus.Fields{
			"host": "host",
		},
	))
}

func GetEntry(ctx context.Context) *logrus.Entry {
	return ctx.Value(loggerKey).(*logrus.Entry)
}
