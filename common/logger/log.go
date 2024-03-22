package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sxp-server/config"
	"sync"
	"time"
)

var Global *zap.SugaredLogger

type ZapLog struct {
	GlobalLog *zap.SugaredLogger
	sync.RWMutex
	fields []string
}

func GetLogger() *ZapLog {
	l := &ZapLog{
		fields: make([]string, 0),
	}
	if Global != nil {
		l.GlobalLog = Global
		return l
	} else {
		IniLogger()
	}
	l.GlobalLog = Global
	return l
}

func IniLogger() {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
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
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl >= logLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= logLevel
	})
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter("./logs/info.log")
	errorWriter := getWriter("./logs/error.log")
	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoLevel), //打印到控制台
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.InfoLevel)) //会显示打日志点的文件名和行数
	Global = log.Sugar()
}

func getWriter(filename string) io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return lumberJackLogger
}

func (l *ZapLog) Debug(args ...interface{}) {
	if len(l.fields) != 0 {
		var s string
		for _, v := range l.fields {
			s += v
		}
		args = append(args, " ", s)
	}
	l.GlobalLog.Debug(args)
}

func (l *ZapLog) Debugf(format string, args ...interface{}) {
	var s string
	if len(l.fields) != 0 {
		for _, v := range l.fields {
			s += v
		}
	}
	l.GlobalLog.Debug(s, " ", fmt.Sprintf(format, args...))
}

func (l *ZapLog) Info(args ...interface{}) {
	if len(l.fields) != 0 {
		var s string
		for _, v := range l.fields {
			s += v
		}
		args = append(args, " ", s)
	}
	l.GlobalLog.Info(args...)
}

func (l *ZapLog) Infof(format string, args ...interface{}) {
	var s string
	if len(l.fields) != 0 {
		for _, v := range l.fields {
			s += v
		}
	}
	l.GlobalLog.Info(s, " ", fmt.Sprintf(format, args...))
}

func (l *ZapLog) Error(args ...interface{}) {
	if len(l.fields) != 0 {
		var s string
		for _, v := range l.fields {
			s += v
		}
		args = append(args, " ", s)
	}
	l.GlobalLog.Error(args...)
}

func (l *ZapLog) Errorf(format string, args ...interface{}) {
	var s string
	if len(l.fields) != 0 {
		for _, v := range l.fields {
			s += v
		}
	}
	l.GlobalLog.Error(s, " ", fmt.Sprintf(format, args...))
}

func (l *ZapLog) Panic(args ...interface{}) {
	if len(l.fields) != 0 {
		var s string
		for _, v := range l.fields {
			s += v
		}
		args = append(args, " ", s)
	}
	l.GlobalLog.Panic(args...)
}

func (l *ZapLog) Panicf(format string, args ...interface{}) {
	var s string
	if len(l.fields) != 0 {
		for _, v := range l.fields {
			s += v
		}
	}
	l.GlobalLog.Panicf(s, " ", fmt.Sprintf(format, args...))
}

func (l *ZapLog) WithFileds(args ...string) *ZapLog {
	l.fields = append(l.fields, args...)
	return l
}
