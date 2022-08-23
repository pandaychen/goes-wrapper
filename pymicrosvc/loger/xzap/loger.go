package xzap

//封装zap，支持requestid打印

import (
	//"github.com/natefinch/lumberjack"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	outWrite     zapcore.WriteSyncer       // IO输出
	debugConsole = zapcore.Lock(os.Stdout) // 控制台标准输出
	once         sync.Once

	defaultGlobalLogger *CtxZapLogger
)

type CtxZapLogger struct {
	//options
	opts     *Option
	filepath string

	zapConfig zap.Config

	logger *zap.Logger
	fields []zapcore.Field

	//for file writer
	writer  zapcore.WriteSyncer
	console zapcore.WriteSyncer
}

// 初始化zap日志配置
func NewCtxZapLoger(optApply ...OptionFunc) {
	once.Do(func() {
		var (
			err          error
			zapLogOption zap.Option
		)
		defaultGlobalLogger = &CtxZapLogger{}

		//初始化日志配置
		defaultGlobalLogger.opts = newZaplogOptions(optApply...)
		defaultGlobalLogger.console = zapcore.Lock(os.Stdout)
		defaultGlobalLogger.filepath = fmt.Sprintf("%s/%s.log", defaultGlobalLogger.opts.LogFileDir, defaultGlobalLogger.opts.AppName)

		if defaultGlobalLogger.opts.Development {
			defaultGlobalLogger.zapConfig = zap.NewDevelopmentConfig()
		} else {
			defaultGlobalLogger.zapConfig = zap.NewProductionConfig()
		}
		//set syncer
		defaultGlobalLogger.SetFileSyncer()
		zapLogOption = defaultGlobalLogger.GetCoresOption()
		if defaultGlobalLogger.logger, err = defaultGlobalLogger.zapConfig.Build(zapLogOption); err != nil {
			panic(err)
		}

		defer defaultGlobalLogger.logger.Sync()
	})
}

func GetCtxZapLoger() *CtxZapLogger {
	if defaultGlobalLogger == nil {
		panic(errors.New("bad vars"))
	}
	return defaultGlobalLogger
}

func GetLoger() *zap.Logger {
	if defaultGlobalLogger == nil {
		panic(errors.New("bad vars"))
	}
	return defaultGlobalLogger.logger
}

//获取带ctx的loger
func (l *CtxZapLogger) GetCtx(ctx context.Context) *zap.Logger {
	loger, ok := ctx.Value(ctxMarkerKey).(*zap.Logger)
	if ok {
		return loger
	}
	return l.logger
}

//向logger增加ctx+日志字段
func (l *CtxZapLogger) WithContextAndFields(ctx context.Context, field ...zap.Field) (context.Context, *zap.Logger) {
	loger := l.logger.With(field...)
	ctx = context.WithValue(ctx, ctxMarkerKey, loger)
	return ctx, loger
}

// 向logger增加ctx
func (l *CtxZapLogger) WithContext(ctx context.Context) *zap.Logger {
	loger, ok := ctx.Value(ctxMarkerKey).(*zap.Logger)
	if ok {
		return loger
	}
	return l.logger
}

//zap底层 API 可以设置缓存，一般主进程退出前使用defer Sync()将缓存同步到文件中
func (l *CtxZapLogger) Sync() error {
	return l.logger.Sync()
}
