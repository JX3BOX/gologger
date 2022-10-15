package rotatefile

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/huyinghuan/lumberjack.v4"
)

var loggerMult map[string]*zap.SugaredLogger

func New(c *lumberjack.Roller, names ...string) *zap.SugaredLogger {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}
	if loggerMult == nil {
		loggerMult = make(map[string]*zap.SugaredLogger)
	}
	if loggerMult[name] != nil {
		panic("logger already init:" + name)
	}
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.InfoLevel)
	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:    "date",
		LevelKey:   "level",
		NameKey:    "name",
		CallerKey:  "caller",
		MessageKey: "msg",
		//FunctionKey:   "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		//EncodeName:   zapcore.FullNameEncoder,
	}
	zl := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(c),
		atom,
	), zap.AddCaller())
	loggerMult[name] = zl.Sugar()
	return zl.Sugar()
}

func Get(names ...string) *zap.SugaredLogger {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}

	if loggerMult[name] != nil {
		return loggerMult[name]
	}
	log.Fatal("logger not found, may be dont init:", name)
	return nil
}
