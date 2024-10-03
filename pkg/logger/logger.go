package logger

import (
	"cli-project/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	Logger *zap.SugaredLogger
)

func init() {
	// Define the custom time format in IST (Indian Standard Time)
	ist := time.FixedZone("IST", 5*3600+1800)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		MessageKey: "msg",
		CallerKey:  "caller",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(ist).Format(time.RFC3339))
		},
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Output paths
	logFilePath := config.LOG_FILE
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Define level enablers
	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel // Console shows Warn and above
	})

	fileLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel // File shows Error and above
	})

	// Set up the core for both console and file logging
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), consoleLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), fileLevel),
	)

	logger := zap.New(core, zap.AddCallerSkip(1), zap.AddCallerSkip(1))
	Logger = logger.Sugar()
}
