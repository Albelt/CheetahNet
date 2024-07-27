package log

import (
	"fmt"
	"go.uber.org/zap"
)

//ILogger的实现,底层使用zap
type logger struct {
	zapLogger *zap.Logger
}

func NewLogger(opts *Options) (*logger, error) {
	if opts == nil {
		return nil, fmt.Errorf("nil Options")
	}

	//level, err := zap.ParseAtomicLevel(opts.Level)
	//if err != nil {
	//	return nil, err
	//}

	config := zap.NewDevelopmentConfig()
	//config.Level = level
	//config.Encoding = opts.Format
	//config.OutputPaths = opts.OutputPaths
	//config.ErrorOutputPaths = opts.ErrOutputPaths

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &logger{zapLogger: zapLogger}, nil
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}
