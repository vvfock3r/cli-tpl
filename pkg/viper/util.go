package viper

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"cli-tpl/pkg/logger"
)

var WatchFunMap = make(map[string]func() error)

// RegisterWatchFunc 注册配置文件监控函数,可以注册多个
// name参数需要唯一,仅用于输出日志时使用
func RegisterWatchFunc(name string, fn func() error) {
	WatchFunMap[name] = fn
}

// StartWatchConfig 开启监控配置文件
func StartWatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// print log
		fileName := e.Name
		fileAbsName, err := filepath.Abs(fileName)
		if err == nil {
			fileName = fileAbsName
		}
		fileName = filepath.ToSlash(fileName)
		operation := strings.ToLower(e.Op.String())
		logger.Warn("config update trigger",
			zap.String("operation", operation),
			zap.String("filename", fileName),
		)

		// Execute all handler functions
		for name, fn := range WatchFunMap {
			err := fn()
			if err == nil {
				logger.Warn("config reload success",
					zap.String("name", name),
					zap.String("detail", "success"),
				)
			} else {
				logger.Warn("config reload ignored",
					zap.String("name", name),
					zap.String("detail", err.Error()),
				)
			}
		}
	})
}
