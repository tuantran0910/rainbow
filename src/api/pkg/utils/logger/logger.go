package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/tuantran0910/rainbow/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ins  *zap.Logger
	once sync.Once
)

func getZapLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s, defaulting to info", level)
	}
}

func getEncoderZap(env string) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	if env == "production" {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func createCoreZap(env string, level zapcore.Level, enableConsole bool) (zapcore.Core, error) {
	// Get the encoder config
	encoder := getEncoderZap(env)

	var outputs []zapcore.WriteSyncer
	if enableConsole || env == "development" {
		outputs = append(outputs, zapcore.AddSync(os.Stdout))
	}

	if env == "production" {
		outputs = append(outputs, zapcore.AddSync(os.Stderr))
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("no log outputs configured")
	}

	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(outputs...),
		zap.NewAtomicLevelAt(level),
	), nil
}

func InitLogger() error {
	var err error
	once.Do(func() {
		// Load the application configurations
		cfg, cfgErr := config.GetConfig()
		if cfgErr != nil {
			err = fmt.Errorf("failed to load the application configurations: %v", cfgErr)
			return
		}

		// Get the log level
		level, levelErr := getZapLevel(cfg.LoggerConfig.LogLevel)
		if levelErr != nil {
			err = fmt.Errorf("failed to get the log level: %v", levelErr)
			return
		}

		// Create the core
		core, coreErr := createCoreZap(cfg.Environment, level, cfg.LoggerConfig.EnableConsole)
		if coreErr != nil {
			err = fmt.Errorf("failed to create the core: %v", coreErr)
			return
		}

		// Create the logger
		var log *zap.Logger
		log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).With(
			zap.String("service", cfg.ServerConfig.ServiceName),
			zap.String("version", cfg.ServerConfig.ServiceVersion),
		)
		if cfg.Environment == "development" {
			log = log.WithOptions(zap.Development(), zap.AddStacktrace(zapcore.ErrorLevel))
		}

		ins = log
	})

	return err
}

func GetLogger() (*zap.Logger, error) {
	if ins == nil {
		if err := InitLogger(); err != nil {
			return nil, err
		}
	}

	return ins, nil
}
