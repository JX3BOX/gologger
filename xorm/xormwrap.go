package xorm

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"xorm.io/xorm/log"
)

// LoggerAdapter wraps a Logger interface as LoggerContext interface
type LoggerAdapter struct {
	logger  *zap.Logger
	level   log.LogLevel
	showSQL bool
}

func XormLogger(logger *zap.Logger) log.ContextLogger {
	return &LoggerAdapter{
		logger: logger,
	}
}

func (l *LoggerAdapter) BeforeSQL(_ log.LogContext) {}

func (l *LoggerAdapter) AfterSQL(lc log.LogContext) {
	if !l.showSQL {
		return
	}
	lg := l.logger

	sql := fmt.Sprintf("%v %v", lc.SQL, lc.Args)

	var level zapcore.Level
	if lc.Err != nil {
		level = zapcore.ErrorLevel
	} else {
		level = zapcore.InfoLevel
	}

	lg.Check(level, "finished sql").Write(zap.String("sql", sql), zap.Duration("execute_time", lc.ExecuteTime), zap.Error(lc.Err))
}

func (l *LoggerAdapter) Debugf(format string, v ...interface{}) {
	l.logger.Sugar().Debugf(format, v...)
}

func (l *LoggerAdapter) Errorf(format string, v ...interface{}) {
	l.logger.Sugar().Errorf(format, v...)
}

func (l *LoggerAdapter) Infof(format string, v ...interface{}) {
	l.logger.Sugar().Infof(format, v...)
}

func (l *LoggerAdapter) Warnf(format string, v ...interface{}) {
	l.logger.Sugar().Warnf(format, v...)
}

func (l *LoggerAdapter) Level() log.LogLevel {
	return l.level
}

func (l *LoggerAdapter) SetLevel(lv log.LogLevel) {
	l.level = lv
}

func (l *LoggerAdapter) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (l *LoggerAdapter) IsShowSQL() bool {
	return l.showSQL
}
