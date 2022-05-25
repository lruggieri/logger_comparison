package iface

import (
	"context"
	"fmt"
	"io"
	"time"
)

type FieldType uint8
type LogLevel uint8

func (l LogLevel) String() string {
	switch l {
	case ILogLevelTrace:
		return "trace"
	case ILogLevelDebug:
		return "debug"
	case ILogLevelInfo:
		return "info"
	case ILogLevelWarning:
		return "warn"
	case ILogLevelError:
		return "error"
	case ILogLevelAlert:
		return "alert"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

type ILogTypedField struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}

type ILogger interface {
	IStdLogger
	IStdTypedLogger

	// SetOutput : sets the logger output
	SetOutput(output io.Writer) (ILogger, error)
	// SetLevel : sets the logging level
	SetLevel(level LogLevel) ILogger
	GetContext() context.Context
}

// IStdLogger : logs using 'message + ...ILogTypedField'. Typed logging is much faster than using interfaces, hence it
// should be the preferred logging way.
// This logging mechanism is structured into fields
type IStdLogger interface {
	Tracef(message string, v ...ILogTypedField)
	Debugf(message string, v ...ILogTypedField)
	Infof(message string, v ...ILogTypedField)
	Warningf(message string, v ...ILogTypedField)
	Errorf(message string, v ...ILogTypedField)
}
type IStdTypedLogger interface {
	String(key string, value string) ILogTypedField
	Bool(key string, value bool) ILogTypedField
	Uint16(key string, value uint16) ILogTypedField
	Uint32(key string, value uint32) ILogTypedField
	Uint64(key string, value uint64) ILogTypedField
	Float32(key string, value float32) ILogTypedField
	Float64(key string, value float64) ILogTypedField
	Int(key string, value int) ILogTypedField
	Int16(key string, value int16) ILogTypedField
	Int32(key string, value int32) ILogTypedField
	Int64(key string, value int64) ILogTypedField
	Duration(key string, value time.Duration) ILogTypedField
	Interface(key string, value interface{}) ILogTypedField
}
type IDefaultTypedLogger struct{}

func (idtl *IDefaultTypedLogger) String(key, value string) ILogTypedField {
	return ILogTypedField{Type: ILogTypeString, Key: key, String: value}
}
func (idtl *IDefaultTypedLogger) Bool(key string, value bool) ILogTypedField {
	var intVal int64 = 0
	if value {
		intVal = 1
	}
	return ILogTypedField{Type: ILogTypeBool, Key: key, Integer: intVal}
}
func (idtl *IDefaultTypedLogger) Uint16(key string, value uint16) ILogTypedField {
	return ILogTypedField{Type: ILogTypeUint16, Key: key, Integer: int64(value)}
}
func (idtl *IDefaultTypedLogger) Uint32(key string, value uint32) ILogTypedField {
	return ILogTypedField{Type: ILogTypeUint32, Key: key, Integer: int64(value)}
}
func (idtl *IDefaultTypedLogger) Uint64(key string, value uint64) ILogTypedField {
	return ILogTypedField{Type: ILogTypeUint64, Key: key, Integer: int64(value)}
}
func (idtl *IDefaultTypedLogger) Float32(key string, value float32) ILogTypedField {
	return ILogTypedField{Type: ILogTypeString, Key: key, String: fmt.Sprintf("%f", value)}
}
func (idtl *IDefaultTypedLogger) Float64(key string, value float64) ILogTypedField {
	return ILogTypedField{Type: ILogTypeString, Key: key, String: fmt.Sprintf("%f", value)}
}
func (idtl *IDefaultTypedLogger) Int(key string, value int) ILogTypedField {
	return ILogTypedField{Type: ILogTypeInt, Key: key, Integer: int64(value)}
}
func (idtl *IDefaultTypedLogger) Int16(key string, value int16) ILogTypedField {
	return ILogTypedField{Type: ILogTypeInt16, Key: key, Integer: int64(value)}
}
func (idtl *IDefaultTypedLogger) Int32(key string, value int32) ILogTypedField {
	return ILogTypedField{Type: ILogTypeInt32, Key: key, Integer: int64(value)}
}
func (idtl *IDefaultTypedLogger) Int64(key string, value int64) ILogTypedField {
	return ILogTypedField{Type: ILogTypeInt64, Key: key, Integer: value}
}

// Duration : always represents a duration in Milliseconds
func (idtl *IDefaultTypedLogger) Duration(key string, value time.Duration) ILogTypedField {
	return ILogTypedField{Type: ILogTypeDuration, Key: key, Integer: int64(value / time.Millisecond)}
}
func (idtl *IDefaultTypedLogger) Interface(key string, value interface{}) ILogTypedField {
	return ILogTypedField{Type: ILogTypeInterface, Key: key, Interface: value}
}
