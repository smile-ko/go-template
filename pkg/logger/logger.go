package logger

import (
	"github/smile-ko/go-template/config"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ILogger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Close()
}

type LoggerZap struct {
	zapLogger *zap.Logger
	syncer    zapcore.WriteSyncer
}

func NewLogger(cfg *config.Config) ILogger {
	encoder := getEncoder(cfg.Log.UseJSON)
	writer := getLogWriter(&cfg.Log, cfg.Log.ConsoleOutput)

	level := getZapLevel(cfg.Log.Level)

	core := zapcore.NewCore(encoder, writer, level)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &LoggerZap{
		zapLogger: logger,
		syncer:    writer,
	}
}

func (l *LoggerZap) Info(msg string, fields ...zap.Field)  { l.zapLogger.Info(msg, fields...) }
func (l *LoggerZap) Debug(msg string, fields ...zap.Field) { l.zapLogger.Debug(msg, fields...) }
func (l *LoggerZap) Warn(msg string, fields ...zap.Field)  { l.zapLogger.Warn(msg, fields...) }
func (l *LoggerZap) Error(msg string, fields ...zap.Field) { l.zapLogger.Error(msg, fields...) }
func (l *LoggerZap) Fatal(msg string, fields ...zap.Field) { l.zapLogger.Fatal(msg, fields...) }

func (l *LoggerZap) Close() {
	_ = l.syncer.Sync()
	_ = l.zapLogger.Sync()
}

func getEncoder(useJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if useJSON {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(cfg *config.Log, consoleOutput bool) zapcore.WriteSyncer {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FileLogName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	if consoleOutput {
		return zapcore.NewMultiWriteSyncer(fileWriter, zapcore.AddSync(os.Stdout))
	}
	return fileWriter
}

func getZapLevel(level string) zapcore.Level {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		log.Printf("invalid log level '%s', fallback to INFO", level)
		lvl = zapcore.InfoLevel
	}
	return lvl
}
