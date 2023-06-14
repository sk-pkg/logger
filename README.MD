# Logger

---
Logger是在zap logger基础上封装的快速使用包，支持日志标准输出、和指定文件输出。

## 安装

```shell
go get -u "github.com/sk-pkg/logger"
```

## 快速开始
默认日志等级为info，默认是以标准输出方式记录日志。
```go
logger, err := New()

if err != nil {
	log.Fatal(err)
}

logger.Sync()

err = errors.New("This is an error testing logger ")

logger = logger.Named("test")

logger.Info("Info")
logger.Debug("Debug")
logger.Error("Error occurs", zap.Error(err))
```
## 可配置项
### 日志驱动
默认驱动为`stdout`为终端标准输出，可选`file`文件驱动，当驱动为`file`时需要同时配置日志路径。
```go
logger, err := New(WithDriver("file"))
if err != nil {
    log.Fatal(err)
}
```
### 日志路径
仅当Driver为file时生效.
```go
logger, err := New(WithDriver("file"),WithLogPath("/data/logs/"))
if err != nil {
    log.Fatal(err)
}
```
### 日志级别
默认日志级别为`info`，可选日志级别为`debug`,`warn`,`error`,`fatal`
```go
logger, err := New(WithLevel("error"))
if err != nil {
    log.Fatal(err)
}
```
### Zap编码配置
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
	EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
}
```
## 参考文档
[Zap原生文档](https://pkg.go.dev/go.uber.org/zap#section-readme)