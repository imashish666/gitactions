package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	Logger *zap.Logger
}

func NewZapLogger() (ZapLogger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewConsoleEncoder(config)
	// fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, _ := os.OpenFile("./logger.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	return ZapLogger{zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))}, nil
}

func (z *ZapLogger) Debug(msg string, fields map[string]interface{}) {
	z.Logger.Debug(msg, getZapFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields map[string]interface{}) {
	z.Logger.Info(msg, getZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields map[string]interface{}) {
	z.Logger.Warn(msg, getZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields map[string]interface{}) {
	z.Logger.Error(msg, getZapFields(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields map[string]interface{}) {
	z.Logger.Fatal(msg, getZapFields(fields)...)
}

func getZapFields(contextMap map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(contextMap))
	for k, v := range contextMap {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}
