package logger

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	DefaultDriver = "stdout"
	DefaultLevel  = zapcore.InfoLevel
)

var DefaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "T",
	LevelKey:       "L",
	NameKey:        "N",
	MessageKey:     "M",
	CallerKey:      "C",
	StacktraceKey:  "S",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

type Option func(*option)

type option struct {
	driver        string                // 日志驱动 stdout, file
	level         zapcore.Level         // 日志级别 debug,info,warn,error,fatal
	logPath       string                // 日志路径，仅当Driver为file时生效
	encoderConfig zapcore.EncoderConfig // Zap编码配置
}

func WithDriver(driver string) Option {
	return func(o *option) {
		o.driver = driver
	}
}

func WithLevel(level string) Option {
	return func(o *option) {
		switch level {
		case "debug":
			o.level = zapcore.DebugLevel
		case "error":
			o.level = zapcore.ErrorLevel
		case "warn":
			o.level = zapcore.WarnLevel
		case "fatal":
			o.level = zapcore.FatalLevel
		default:
			o.level = zapcore.InfoLevel
		}
	}
}

func WithLogPath(path string) Option {
	return func(o *option) {
		o.logPath = path
	}
}

func WithEncoderConfig(config zapcore.EncoderConfig) Option {
	return func(o *option) {
		o.encoderConfig = config
	}
}

func New(opts ...Option) (*zap.Logger, error) {
	opt := &option{driver: DefaultDriver, level: DefaultLevel, encoderConfig: DefaultEncoderConfig}
	for _, f := range opts {
		f(opt)
	}

	jsonEncoder := zapcore.NewJSONEncoder(opt.encoderConfig)

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	// 标准输出
	if opt.driver == "stdout" {
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	}

	// 日志文件输出
	if opt.driver == "file" {
		// 例子：/data/logs/logger/2021-05-17.log
		hook, err := rotatelogs.New(opt.logPath + "%Y-%m-%d.log")

		if err != nil {
			return nil, err
		}

		core = zapcore.NewTee(core,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(hook),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	logger := zap.New(core,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	return logger, nil
}
