package xzap

import (
	"fmt"
	"path/filepath"
)

type OptionFunc func(*Option)

type Option struct {
	Tags []string //TODO：记录zap的fileds

	Development        bool
	LogFileDir         string
	AppName            string
	MaxSize            int
	MaxBackups         int
	MaxAge             int
	Level              string
	EnableCompress     bool
	EnableWriteFile    bool
	EnableWriteConsole bool
}

// 初始化log的配置
func newZaplogOptions(optFunctions ...OptionFunc) *Option {
	var (
		curPath string
		err     error
	)
	option := &Option{
		Development:        true,
		AppName:            "test-apps",
		MaxSize:            100,
		MaxBackups:         60,
		MaxAge:             30,
		Level:              "debug",
		EnableCompress:     false,
		EnableWriteFile:    false,
		EnableWriteConsole: true,
	}

	if curPath, err = filepath.Abs(filepath.Dir(filepath.Join("."))); err != nil {
		panic(err)
	}

	option.LogFileDir = fmt.Sprintf("%s/log/", curPath)

	for _, o := range optFunctions {
		//apply options
		o(option)
	}
	return option
}

func SetDevelopment(development bool) OptionFunc {
	return func(options *Option) {
		options.Development = development
	}
}

func SetLogFileDir(logFileDir string) OptionFunc {
	return func(options *Option) {
		options.LogFileDir = logFileDir
	}
}

func SetAppName(appName string) OptionFunc {
	return func(options *Option) {
		options.AppName = appName
	}
}

func SetMaxSize(maxSize int) OptionFunc {
	return func(options *Option) {
		options.MaxSize = maxSize
	}
}
func SetMaxBackups(maxBackups int) OptionFunc {
	return func(options *Option) {
		options.MaxBackups = maxBackups
	}
}
func SetMaxAge(maxAge int) OptionFunc {
	return func(options *Option) {
		options.MaxAge = maxAge
	}
}

func SetLevel(level string) OptionFunc {
	return func(options *Option) {
		options.Level = level
	}
}

func SetWriteFile(writeFile bool) OptionFunc {
	return func(options *Option) {
		options.EnableWriteFile = writeFile
	}
}

func SetLogCompress(isOn bool) OptionFunc {
	return func(options *Option) {
		options.EnableCompress = isOn
	}
}

func SetWriteConsole(writeConsole bool) OptionFunc {
	return func(options *Option) {
		options.EnableWriteConsole = writeConsole
	}
}
