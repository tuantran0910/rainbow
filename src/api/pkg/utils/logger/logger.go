package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ins  *zap.Logger
	once sync.Once
)

func getEncoderZap() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			LevelKey:     "level",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeCaller: zapcore.ShortCallerEncoder,
			CallerKey:    "caller",
			EncodeLevel:  customEncodeLevel,
		},
	)
}

func customEncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + l.CapitalString() + "]")
}

func getLogWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func NewLogger() *zap.Logger {
	encoder := getEncoderZap()
	sync := getLogWriter()
	core := zapcore.NewCore(encoder, sync, zapcore.DebugLevel)
	logger := zap.New(core)
	return logger
}

func GetLogger() *zap.Logger {
	once.Do(func() {
		ins = NewLogger()
	})

	return ins
}
