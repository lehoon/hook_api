package logger

import (
	"github.com/lehoon/hook_api/v2/library/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.SugaredLogger
)

func init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	atomLevel := zap.NewAtomicLevelAt(zapcore.Level(config.GetLoggerLevel()))

	loggerHook := lumberjack.Logger{
		Filename:   config.GetLoggerPath(),
		MaxSize:    config.GetLoggerMaxSize(),
		MaxAge:     config.GetLoggerMaxAge(),
		MaxBackups: config.GetLoggerMaxBackup(),
		Compress:   config.GetLoggerCompress(),
	}

	loggerCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&loggerHook)),
		atomLevel,
	)

	caller := zap.AddCaller()
	development := zap.Development()
	field := zap.Fields(zap.String("serviceName", "hook_api"))
	logger := zap.New(loggerCore, caller, development, field)
	logger.Info("logger初始化完成")
	log = logger.Sugar()
}

func Debug(message string) {
	log.Debug(message)
}

func Info(message string) {
	log.Info(message)
}

func Error(message string) {
	log.Error(message)
}

func Log() *zap.SugaredLogger {
	return log
}

func Sync() {
	log.Sync()
}
