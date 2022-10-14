package rotatefile

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/huyinghuan/lumberjack.v4"
)

var loggerMult map[string]*zap.SugaredLogger

func New(name string, c *lumberjack.Roller) *zap.SugaredLogger {
	if name == "" {
		name = "default"
	}
	if loggerMult == nil {
		loggerMult = make(map[string]*zap.SugaredLogger)
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
	name := ""
	if len(names) > 0 {
		name = names[0]
	} else {
		for k, _ := range loggerMult {
			name = k
			break
		}
	}

	if loggerMult[name] != nil {
		return loggerMult[name]
	}
	log.Fatal("logger not found, may be dont init:", name)
	return nil
}
