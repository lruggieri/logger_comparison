package logger

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"logging/src/iface"

	"context"
	"fmt"
	"io"
	"net/url"
)

const (
	zapCustomWriterBaseScheme = "cwriter"
)

var (
	zapLogLevelMapping = map[iface.LogLevel]zapcore.Level{
		iface.ILogLevelTrace:   zapcore.DebugLevel,
		iface.ILogLevelDebug:   zapcore.DebugLevel,
		iface.ILogLevelInfo:    zapcore.InfoLevel,
		iface.ILogLevelWarning: zapcore.DebugLevel,
		iface.ILogLevelError:   zapcore.ErrorLevel,
	}
	zapCommonOptions = []zap.Option{
		zap.AddCallerSkip(1), // to skip this zap.go call and the flood_prevention wrapper function
	}
)

type ZapLogger struct {
	iface.IDefaultTypedLogger

	ctx         context.Context
	config      zap.Config
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger

	additionalCallerSkip int
	isProduction         bool
}

func NewZapLogger(
	iCtx context.Context, isProduction bool, additionalCallerSkip int) (logger iface.ILogger, err error) {

	var zapLogger *zap.Logger
	var config zap.Config
	if isProduction {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// default level should always be INFO
	config.Level.SetLevel(zap.InfoLevel)

	opts := zapCommonOptions
	opts = append(opts, zap.AddCallerSkip(additionalCallerSkip))
	zapLogger, err = config.Build(opts...)
	if err != nil {
		return nil, err
	}

	logger = &ZapLogger{
		ctx:                  iCtx,
		additionalCallerSkip: additionalCallerSkip,
		isProduction:         isProduction,
		config:               config,
		logger:               zapLogger,
		sugarLogger:          zapLogger.Sugar(),
	}

	return logger, nil
}

func (zl *ZapLogger) GetContext() context.Context {
	return zl.ctx
}

func getZapType(t iface.FieldType) zapcore.FieldType {
	switch t {
	case iface.ILogTypeString:
		return zapcore.StringType
	case iface.ILogTypeBool:
		return zapcore.BoolType
	case iface.ILogTypeUint16:
		return zapcore.Uint16Type
	case iface.ILogTypeUint32:
		return zapcore.Uint32Type
	case iface.ILogTypeUint64:
		return zapcore.Uint64Type
	case iface.ILogTypeInt16:
		return zapcore.Int16Type
	case iface.ILogTypeInt32:
		return zapcore.Int32Type
	case iface.ILogTypeInt, iface.ILogTypeInt64:
		return zapcore.Int64Type
	case iface.ILogTypeDuration:
		return zapcore.Int64Type // avoid third party libraries to do further translation on complex types
	case iface.ILogTypeInterface:
		return zapcore.ReflectType
	default:
		return zapcore.ReflectType
	}
}
func getZapFields(fields []iface.ILogTypedField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {

		zapFields = append(zapFields, zap.Field{
			Key:       f.Key,
			Type:      getZapType(f.Type),
			String:    f.String,
			Integer:   f.Integer,
			Interface: f.Interface,
		})
	}
	return zapFields
}

// deduplicateZapFields: remove duplicated zap fields by taking the first field found for each key
func deduplicateZapFields(fields []zap.Field) []zap.Field {
	seenKeys := make(map[string]struct{})
	finalFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		if _, ok := seenKeys[f.Key]; !ok {
			seenKeys[f.Key] = struct{}{}
			finalFields = append(finalFields, f)
		}
	}
	return finalFields
}
func (zl *ZapLogger) getLogLevel() iface.LogLevel {
	currentZapLevel := zl.config.Level.Level()
	for level, zapLevel := range zapLogLevelMapping {
		if currentZapLevel == zapLevel {
			return level
		}
	}
	return iface.ILogLevelInfo // default fallback
}
func (zl *ZapLogger) SetSkipLineEnding(iValue bool) {
	zl.config.EncoderConfig.SkipLineEnding = iValue
}

func (zl *ZapLogger) Tracef(message string, v ...iface.ILogTypedField) {
	zl.logger.Debug(message, getZapFields(v)...)
}
func (zl *ZapLogger) Debugf(message string, v ...iface.ILogTypedField) {
	zl.logger.Debug(message, getZapFields(v)...)
}
func (zl *ZapLogger) Infof(message string, v ...iface.ILogTypedField) {
	zl.logger.Info(message, getZapFields(v)...)
}
func (zl *ZapLogger) Warningf(message string, v ...iface.ILogTypedField) {
	zl.logger.Warn(message, getZapFields(v)...)
}
func (zl *ZapLogger) Errorf(message string, v ...iface.ILogTypedField) {
	zl.logger.Error(message, getZapFields(v)...)
}

type customWriter struct {
	io.Writer
}

func (cw customWriter) Close() error {
	return nil
}
func (cw customWriter) Sync() error {
	return nil
}
func (zl *ZapLogger) SetOutput(output io.Writer) (logger iface.ILogger, err error) {
	cw := customWriter{output}
	newUniqueScheme := fmt.Sprintf("%s%s", zapCustomWriterBaseScheme, uuid.New().String())
	err = zap.RegisterSink(newUniqueScheme, func(u *url.URL) (zap.Sink, error) {
		return cw, nil
	})
	if err != nil {
		return zl, err
	}

	zl.config.OutputPaths = []string{newUniqueScheme + ":something"}
	oldLogger := zl.logger
	opts := zapCommonOptions
	opts = append(opts, zap.AddCallerSkip(zl.additionalCallerSkip))
	zl.logger, err = zl.config.Build(opts...)
	if err != nil {
		return zl, err
	}
	_ = oldLogger.Sync() // flush old logger
	return zl, nil
}
func (zl *ZapLogger) SetLevel(level iface.LogLevel) iface.ILogger {
	var loggerLevel zapcore.Level
	if zapLevel, ok := zapLogLevelMapping[level]; ok {
		loggerLevel = zapLevel
	} else {
		loggerLevel = zapcore.InfoLevel // default fallback
	}

	zl.config.Level.SetLevel(loggerLevel)
	return zl
}
