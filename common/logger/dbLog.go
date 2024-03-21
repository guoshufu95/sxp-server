package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"sxp-server/config"
	"time"
)

var (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s [%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

type gormLogger struct {
	ZapLogger                           *ZapLog
	LogLevel                            logger.LogLevel
	SlowThreshold                       time.Duration
	SkipCallerLookup                    bool
	IgnoreRecordNotFoundError           bool
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {

		l.ZapLogger.Infof(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Infof(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > logger.Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				l.ZapLogger.Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.ZapLogger.Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.ZapLogger.Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.ZapLogger.Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel == logger.Info:
			sql, rows := fc()
			if rows == -1 {
				l.ZapLogger.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.ZapLogger.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

func NewGormLogger() logger.Interface {
	logLevel := zap.DebugLevel
	switch config.Conf.Logger.Level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.InfoLevel
	}
	return &gormLogger{
		ZapLogger:                 GetLogger(),
		LogLevel:                  logger.LogLevel(logLevel),
		SlowThreshold:             200 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		infoStr:                   infoStr,
		warnStr:                   warnStr,
		errStr:                    errStr,
		traceStr:                  traceStr,
		traceWarnStr:              traceWarnStr,
		traceErrStr:               traceErrStr,
	}
}
