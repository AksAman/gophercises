package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitializeLogger() {

	core := zapcore.NewTee(
		zapcore.NewCore(getConsoleEncoder(), zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(getJSONEncoder(), getLogWriter(), zapcore.DebugLevel),
	)
	Logger = zap.New(core, zap.AddCaller()).Sugar()
	defer Logger.Sync()

	Logger.Info("Logger initialized")
}

func getEncoderConfig() zapcore.EncoderConfig {
	baseConfig := zap.NewProductionEncoderConfig()
	baseConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	baseConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	baseConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return baseConfig
}

func getJSONEncoder() zapcore.Encoder {
	config := getEncoderConfig()
	return zapcore.NewJSONEncoder(config)
}

func getConsoleEncoder() zapcore.Encoder {
	config := getEncoderConfig()
	return zapcore.NewConsoleEncoder(config)
}

func getLogWriter() zapcore.WriteSyncer {
	logsPath := "logs"
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		os.Mkdir(logsPath, os.ModePerm)
	}
	file, _ := os.Create(filepath.Join(logsPath, "sitemap.log"))
	return zapcore.AddSync(file)
}
