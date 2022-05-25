package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

type WrapperLogrus struct {
	*logrus.Entry
}

func NewWrapperLogrus() *WrapperLogrus {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	l := logrus.StandardLogger()
	l.SetReportCaller(true)
	return &WrapperLogrus{
		Entry: l.WithContext(context.Background()),
	}
}

func (logger *WrapperLogrus) SetLevel(level string) {
	if lvl, err := logrus.ParseLevel(level); err == nil {
		logger.Logger.SetLevel(lvl)
		return
	}
	logger.Logger.SetLevel(logrus.InfoLevel)
}

func (logger *WrapperLogrus) WithField(key string, value interface{}) *WrapperLogrus {
	return &WrapperLogrus{
		Entry: logger.Entry.WithFields(logrus.Fields{key: value}),
	}
}

func (logger *WrapperLogrus) Info(ctx context.Context, args ...interface{}) {
	logger.log(ctx, logrus.InfoLevel, args...)
}

func (logger *WrapperLogrus) Debug(ctx context.Context, args ...interface{}) {
	logger.log(ctx, logrus.DebugLevel, args...)
}

func (logger *WrapperLogrus) Error(ctx context.Context, args ...interface{}) {
	logger.log(ctx, logrus.ErrorLevel, args...)
}

func (logger *WrapperLogrus) Warn(ctx context.Context, args ...interface{}) {
	logger.log(ctx, logrus.WarnLevel, args...)
}

func (logger *WrapperLogrus) log(ctx context.Context, level logrus.Level, args ...interface{}) {
	if !logger.Logger.IsLevelEnabled(level) {
		return
	}
	entry := logger.Entry
	entry.Log(level, args...)
}

func (logger *WrapperLogrus) Infof(ctx context.Context, format string, args ...interface{}) {
	logger.logf(ctx, logrus.InfoLevel, format, args...)
}

func (logger *WrapperLogrus) Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.logf(ctx, logrus.DebugLevel, format, args...)
}

func (logger *WrapperLogrus) Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.logf(ctx, logrus.ErrorLevel, format, args...)
}

func (logger *WrapperLogrus) Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.logf(ctx, logrus.WarnLevel, format, args...)
}

func (logger *WrapperLogrus) logf(ctx context.Context, level logrus.Level, format string, args ...interface{}) {
	if !logger.Logger.IsLevelEnabled(level) {
		return
	}
	entry := logger.Entry

	entry.Logf(level, format, args...)
}
