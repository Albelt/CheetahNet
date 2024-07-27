package log

const (
	OUTPUT_STD = "stdout"
	ERR_STD    = "stderr"

	FORMAT_CONSOLE = "console"
	FORMAT_JSON    = "json"

	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
)

//日志配置
type Options struct {
	OutputPaths    []string
	ErrOutputPaths []string
	Level          string
	Format         string
}

type Option func(o *Options)

func NewOptions(opts ...Option) *Options {
	//默认配置
	options := &Options{
		OutputPaths:    []string{OUTPUT_STD},
		ErrOutputPaths: []string{ERR_STD},
		Level:          LOG_LEVEL_INFO,
		Format:         FORMAT_CONSOLE,
	}

	//用户配置
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithFormat(format string) Option {
	return func(o *Options) {
		o.Format = format
	}
}