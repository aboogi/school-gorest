package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(zap.InfoLevel.String())
	if err != nil {
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	encconf := zap.NewProductionEncoderConfig()
	encconf.TimeKey = "@timestamp"
	encconf.EncodeTime = zapcore.RFC3339TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encconf)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), lvl)

	core := zapcore.NewTee(consoleCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}
