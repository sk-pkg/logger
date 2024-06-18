package logger

import (
	"context"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	DefaultDriver = "stdout"
	DefaultLevel  = zapcore.InfoLevel
	TraceIDKey    = "trace_id"
)

type (
	Option func(*option)

	option struct {
		driver        string                // 日志驱动 stdout, file
		level         zapcore.Level         // 日志级别 debug,info,warn,error,fatal
		logPath       string                // 日志路径，仅当Driver为file时生效
		encoderConfig zapcore.EncoderConfig // Zap编码配置
	}

	Manager struct {
		Zap *zap.Logger
	}
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

func New(opts ...Option) (*Manager, error) {
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
		zap.AddCallerSkip(1), // 跳过封装函数这一层
		zap.ErrorOutput(stderr),
	)

	return &Manager{Zap: logger}, nil
}

// getTraceIDFromContext 从上下文中提取 TraceID
func getTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

func (m *Manager) getLoggerWithTraceID(ctx context.Context) *zap.Logger {
	traceID := getTraceIDFromContext(ctx)
	if traceID != "" {
		return m.Zap.With(zap.String("TraceID", traceID))
	}
	return m.Zap
}

func (m *Manager) Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Info(msg, fields...)
}

func (m *Manager) Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Error(msg, fields...)
}

func (m *Manager) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Debug(msg, fields...)
}

func (m *Manager) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Warn(msg, fields...)
}

func (m *Manager) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Fatal(msg, fields...)
}

func (m *Manager) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Panic(msg, fields...)
}

func (m *Manager) Sync() {
	m.Zap.Sync()
}

func (m *Manager) Named(name string) *zap.Logger {
	return m.Zap.Named(name)
}
