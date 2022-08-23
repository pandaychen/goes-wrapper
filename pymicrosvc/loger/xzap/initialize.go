package xzap

import (
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//	设置zap.AddSync
func (l *CtxZapLogger) SetFileSyncer() {
	lumberhooker := &lumberjack.Logger{
		Filename:   l.filepath,
		MaxSize:    l.opts.MaxSize,
		MaxBackups: l.opts.MaxBackups,
		MaxAge:     l.opts.MaxAge,
		Compress:   l.opts.EnableCompress,
		LocalTime:  true,
	}

	l.writer = zapcore.AddSync(lumberhooker)
}

// 设置启用日志类型
func (l *CtxZapLogger) GetCoresOption() zap.Option {
	var (
		cores []zapcore.Core
	)
	//use json encoder
	json_encoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)

	//设置日志等级下限
	priority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= l.GetLogLevel()
	})

	if l.opts.EnableWriteFile && l.writer != nil {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(json_encoder, l.writer, priority),
		}...)
	}
	if l.opts.EnableWriteConsole && l.console != nil {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(json_encoder, l.console, priority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

//从配置文件中获取日志配置等级
func (l *CtxZapLogger) GetLogLevel() (level zapcore.Level) {
	var (
		logLvl = strings.ToLower(l.opts.Level)
	)
	switch logLvl {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		//use default
		return zapcore.DebugLevel
	}
}
