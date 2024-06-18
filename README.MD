
# Logger
Logger 是在 zap logger 基础上封装的快速使用包。
## 特性
- 支持日志标准输出、和指定文件输出。
- 支持TraceID。


## 安装
```shell
go get -u "github.com/sk-pkg/logger"
````

## 快速开始
默认日志等级为 info，默认是以标准输出方式记录日志。
```go
package main

import (
    "context"
    "errors"
    "github.com/sk-pkg/logger"
    "go.uber.org/zap"
    "log"
)

func main() {
    // 创建 Logger 管理器实例并设置为 debug 级别
    loggerManager, err := logger.New(logger.WithLevel("debug"))
    if err != nil {
        log.Fatal(err)
    }

    // 设置上下文中的 TraceID
    ctx := context.WithValue(context.Background(), logger.TraceIDKey, "123456")

    defer loggerManager.Sync()

    // 记录不同级别的日志
    loggerManager.Debug(ctx, "Debug")
    loggerManager.Info(ctx, "Info", zap.String("info", "value"))
    loggerManager.Error(ctx, "error", zap.Error(errors.New("debug info")))
}
```
## 可配置项
### 日志驱动
默认驱动为 `stdout`，即终端标准输出。可选 `file` 文件驱动，当驱动为 `file` 时需要同时配置日志路径。
```go
loggerManager, err := logger.New(logger.WithDriver("file"))
if err != nil {
    log.Fatal(err)
}
```
### 日志路径
仅当 Driver 为 file 时生效。
```go
loggerManager, err := logger.New(logger.WithDriver("file"), logger.WithLogPath("/data/logs/"))
if err != nil {
    log.Fatal(err)
}
```
### 日志级别
默认日志级别为 `info`，可选日志级别为 `debug`, `warn`, `error`, `fatal`。
```go
loggerManager, err := logger.New(logger.WithLevel("fatal"))
if err != nil {
    log.Fatal(err)
}
```
### Zap 编码配置
可以自定义 zap 的编码配置。
```go
var DefaultEncoderConfig = zapcore.EncoderConfig{
    TimeKey:        "time",
    LevelKey:       "level",
    NameKey:        "namespace",
    MessageKey:     "msg",
    CallerKey:      "caller",
    StacktraceKey:  "stacktrace",
    LineEnding:     zapcore.DefaultLineEnding,
    EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
    EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
    EncodeDuration: zapcore.SecondsDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,    // 短路径编码器
}
```
## 代码示例
### 创建 Logger 管理器

```go
package main

import (
	"context"
	"github.com/sk-pkg/logger"
	"log"
)

func main() {
	// 创建 Logger 管理器实例
	loggerManager, err := logger.New(
		logger.WithDriver("stdout"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer loggerManager.Sync()

	// 设置上下文中的 TraceID
	ctx := context.WithValue(context.Background(), logger.TraceIDKey, "test-trace-id")

	// 使用 Logger 记录日志
	loggerManager.Info(ctx, "This is an info message")
	loggerManager.Error(ctx, "This is an error message")
	loggerManager.Debug(ctx, "This is a debug message")
	loggerManager.Warn(ctx, "This is a warn message")
}
```
### 自定义编码配置

```go
package main

import (
	"context"
	"github.com/sk-pkg/logger"
	"go.uber.org/zap/zapcore"
	"log"
)

func main() {
	customEncoderConfig := zapcore.EncoderConfig{
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

	loggerManager, err := logger.New(
		logger.WithEncoderConfig(customEncoderConfig),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer loggerManager.Sync()

	ctx := context.Background()

	// 使用自定义配置的 Logger 记录日志
	loggerManager.Info(ctx, "This is a custom configured info message")
	loggerManager.Error(ctx, "This is a custom configured error message")
}
```
### 日志文件输出

```go
package main

import (
	"context"
	"github.com/sk-pkg/logger"
	"log"
)

func main() {
	loggerManager, err := logger.New(
		logger.WithDriver("file"),
		logger.WithLogPath("/data/logs/"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer loggerManager.Sync()

	ctx := context.Background()

	// 使用记录到文件的 Logger 记录日志
	loggerManager.Info(ctx, "This is an info message written to file")
	loggerManager.Error(ctx, "This is an error message written to file")
}
```
### 获取带有 TraceID 的 Logger
```go
package main

import (
    "context"
    "github.com/sk-pkg/logger"
    "go.uber.org/zap"
    "log"
)

func main() {
    loggerManager, err := logger.New()
    if err != nil {
        log.Fatal(err)
    }
    defer loggerManager.Sync()

    // 设置上下文中的 TraceID
    ctx := context.WithValue(context.Background(), logger.TraceIDKey, "test-trace-id")

    // 记录带有 TraceID 的日志
    loggerManager.Info(ctx, "This is an info message with TraceID", zap.String("additionalField", "additionalValue"))
}
```
## 参考文档
[Zap 原生文档](https://pkg.go.dev/go.uber.org/zap#section-readme)