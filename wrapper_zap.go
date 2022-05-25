package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sirupsen/logrus"
)

type WrapperZap struct {
	*zap.Logger

	config zap.Config
}

func NewWrapperZap() *WrapperZap {
	zapConfig := zap.NewProductionConfig()
	opts := []zap.Option{
		zap.AddCallerSkip(0),
	}
	zapConfig.Level.SetLevel(zap.InfoLevel)
	l, _ := zapConfig.Build(opts...)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	return &WrapperZap{
		Logger: l,
		config: zapConfig,
	}
}

func (logger *WrapperZap) SetLevel(level zapcore.Level) {
	logger.config.Level.SetLevel(level)
}
