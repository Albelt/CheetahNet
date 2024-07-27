package log

import (
	"EagleNet/eiface"
	"sync"
)

//全局默认ILogger
var (
	defaultLogger eiface.ILogger
	mu            sync.Mutex
)

//初始化defaultLogger
func init() {
	mu.Lock()
	defer mu.Unlock()

	var err error
	if defaultLogger == nil {
		defaultLogger, err = NewLogger(NewOptions())
		if err != nil {
			panic(err)
		}
	}
}

// ----------------- 提供外部可以直接使用的日志函数 ---------------------

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}