package gologger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLogger *zap.Logger
var loggerMult map[string]*zap.SugaredLogger
var atomicLevel zap.AtomicLevel
var sugar *zap.SugaredLogger

func GetSingleInstance() *zap.Logger {
	return zapLogger
}

func InitLogger(l Level) {
	level := zapcore.Level(l)
	if zapLogger != nil {
		if level != atomicLevel.Level() {
			Infof("zap logger change level from %s to %s", atomicLevel.Level(), level)
			atomicLevel.SetLevel(level)
		}
		return
	}
	atomicLevel = zap.NewAtomicLevel()
	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:    "date",
		LevelKey:   "level",
		NameKey:    "name",
		CallerKey:  "caller",
		MessageKey: "msg",
		// FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		//	EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		//EncodeName:   zapcore.FullNameEncoder,
	}
	atomicLevel.SetLevel(level)
	if level == zapcore.DebugLevel {
		encoderCfg.FunctionKey = "func"
	}
	zapLogger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	), zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = zapLogger.Sugar()
}

func getSugarLogger() *zap.SugaredLogger {
	if zapLogger == nil {
		InitLogger(InfoLevel)
	}
	return sugar
}

func Debug(args ...interface{}) {
	getSugarLogger().Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	getSugarLogger().Debugf(template, args...)
}

func Info(args ...interface{}) {
	getSugarLogger().Info(args...)
}

func Infof(template string, args ...interface{}) {
	getSugarLogger().Infof(template, args...)
}

func Infow(template string, args ...interface{}) {
	getSugarLogger().Infow(template, args...)
}

func Warn(args ...interface{}) {
	getSugarLogger().Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	getSugarLogger().Warnf(template, args...)
}

func Error(args ...interface{}) {
	getSugarLogger().Error(args...)
}

func Errorf(template string, args ...interface{}) {
	getSugarLogger().Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	getSugarLogger().DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	getSugarLogger().DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	getSugarLogger().Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	getSugarLogger().Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	getSugarLogger().Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	getSugarLogger().Fatalf(template, args...)
}
