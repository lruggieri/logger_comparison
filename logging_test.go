package logging

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"io"
	"logging/src/logger"
	"net/url"
	"testing"
)

func Benchmark_Zap(b *testing.B) {
	// setting output
	var bts bytes.Buffer
	bWriter := bufio.NewWriter(&bts)
	cw := customWriter{bWriter}

	// preparing Zap configurations and initializing logger
	newUniqueScheme := fmt.Sprintf("%s%s", "cwriter", uuid.New().String())
	_ = zap.RegisterSink(newUniqueScheme, func(u *url.URL) (zap.Sink, error) {
		return cw, nil
	})
	pc := zap.NewProductionConfig()
	pc.OutputPaths = []string{newUniqueScheme + ":something"}
	zl, _ := pc.Build()

	// testing
	for n := 0; n < b.N; n++ {
		zl.Info("ciao",
			zap.String("string1", "value1"),
			zap.Int("int1", 1),
			zap.Bool("bool1", true),
		)
	}
}
func Benchmark_Logrus(b *testing.B) {
	// setting output
	var bts bytes.Buffer
	bWriter := bufio.NewWriter(&bts)
	cw := customWriter{bWriter}

	// preparing Zap configurations and initializing logger
	lg := logrus.New()
	lg.SetOutput(cw)

	for n := 0; n < b.N; n++ {
		lg.WithFields(logrus.Fields{
			"string1": "value1",
			"int1":    1,
			"bool1":   true,
		}).Info("ciao")
	}
}

func Test_LogrusCallerWithWrapper(t *testing.T) {
	l := NewWrapperLogrus()
	l.Info(context.Background(), "This is an INFO log")
}
func Test_ZapCallerWithWrapper(t *testing.T) {
	l := NewWrapperZap()

	l.Info("This is an INFO log",
		zap.String("string_field", "I am a string"),
		zap.Int("int_field", 42),
		zap.Bool("bool_field", true),
		zap.Reflect("interface_field", struct{ A int }{123}),
	)
}
func Test_ZapCallerWithInterface(t *testing.T) {
	l, _ := logger.NewZapLogger(context.Background(), true, 0)

	l.Infof("This is an INFO log",
		l.String("string_field", "I am a string"),
		l.Int("int_field", 42),
		l.Bool("bool_field", true),
		l.Interface("interface_field", struct{ A int }{123}),
	)
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
