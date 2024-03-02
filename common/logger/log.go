package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
	"time"
)

var Global *zap.SugaredLogger

type ZapLog struct {
	GlobalLog *zap.SugaredLogger
	sync.RWMutex
	fields map[string]interface{}
}

func GetLogger() *ZapLog {
	l := &ZapLog{
		fields: make(map[string]interface{}),
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
	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
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

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel)) //会显示打日志点的文件名和行数
	Global = log.Sugar()
}

func getWriter(filename string) io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return lumberJackLogger
}

func (l *ZapLog) Info(args ...interface{}) {
	l.GlobalLog.Info(args)
}

func (l *ZapLog) Infof(format string, args ...interface{}) {
	l.GlobalLog.Infof(format, args...)
}

func (l *ZapLog) Error(args ...interface{}) {
	l.GlobalLog.Error(args...)
}

func (l *ZapLog) Errorf(format string, args ...interface{}) {
	l.GlobalLog.Errorf(format, args...)
}

func (l *ZapLog) Panic(args ...interface{}) {
	l.GlobalLog.Panic(args...)
}

func (l *ZapLog) Panicf(format string, args ...interface{}) {
	l.GlobalLog.WithOptions()
	l.GlobalLog.Panicf(format, args...)
}
