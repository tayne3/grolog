<p align="center" style="font-size: 24px; font-weight: bold; color: #797979;">Grolog</p>

**English** | [中文](README_zh.md)

Grolog 是一个轻量级的 Go 日志库,旨在提供简单、高性能和可配置的日志记录功能。它支持多种日志级别、日志格式以及异步和同步模式。此外, Grolog 还提供了灵活的配置选项,使您可以根据需求自定义日志行为。

## 特性

- **多种日志级别**: 变量(VerBose)、调试(Debug)、跟踪(Trace)、警告(Warning)、错误(Error)和致命错误(Fatal)。
- **多种日志格式**: 基本(Basic)、简要(Brief)和详细(Detail)。
- **异步和同步模式**: 支持异步和同步两种日志记录模式,可根据需求选择。
- **文件日志记录**: 支持将日志写入文件,可配置文件大小限制、文件数量限制和过期时间。
- **控制台日志记录**: 支持将日志输出到控制台。
- **高度可配置**: 提供多种配置选项,可根据需求自定义日志行为。
- **JSON 配置支持**: 支持从 JSON 配置文件加载日志配置。

## 快速入门

### 基本用法

```go
import (
    "grolog"
)

func main() {
    // 创建日志器
    logger := grolog.Default() // 等同于 `logger := grolog.New(nil)`
	defer logger.Close()

    // 记录日志
    logger.Debug("This is a debug message.\n")
    logger.Warningf("Warning: %s\n", "Something went wrong.")
    logger.Errorln("An error occurred.")
}
```

### 配置选项

grolog 提供了多种配置选项,可以根据需求自定义日志行为。以下是一些常用的配置选项:

- `WithFatalHandling(handling func(*Logger, any))`: 设置异常日志处理函数,用于在发生致命错误时执行特定操作。
- `WithMsgCallback(handler func(int, string))`: 设置日志消息回调函数,用于在记录日志时执行自定义操作。
- `WithGoExec(exec func(f func()))`: 设置异步执行函数,用于支持外部 goroutine 池。
- `WithLevel(level int)`: 设置日志级别,可选值为 `LevelVerBose`、`LevelDebug`、`LevelTrace`、`LevelWarning`、`LevelError` 和 `LevelFatal`。
- `WithStyle(style int)`: 设置日志格式,可选值为 `StyleBasic`、`StyleBrief` 和 `StyleDetail`。
- `WithEnableAsyn(asyn bool)`: 启用或禁用异步日志记录模式 (同步模式的性能可能会优于异步模式，但异步模式下资源使用更加可控)。
- `WithEnableFileTime(enable bool)`: 启用或禁用包含时间信息的日志文件名。
- `WithDisableSave(save bool)`: 禁用或启用文件日志记录。
- `WithDisablePrint(print bool)`: 禁用或启用控制台日志记录。
- `WithAsynMaxGor(max int)`: 设置异步执行数量上限,启用异步模式时有效。
- `WithAsynMaxBuffer(max int)`: 设置异步消息缓冲大小,启用异步模式时有效。
- `WithWriteBufferSize(size int)`: 设置日志文件缓冲大小,启用日志文件时有效。
- `WithFileDir(dir string)`: 设置日志文件目录,启用日志文件时有效。
- `WithFileName(name string)`: 设置日志文件名称,启用日志文件时有效。
- `WithMaxFileSize(maxSize int64)`: 设置单个日志文件的最大大小,启用日志文件时有效。
- `WithMaxFileCount(maxCount int)`: 设置最大日志文件数量,启用日志文件时有效。
- `WithFlashInterval(interval string)`: 设置日志刷新间隔,启用日志文件时有效。
- `WithExpireTime(expire string)`: 设置日志文件过期时间,启用日志文件时有效。

示例:

```go
import (
    "grolog"
)

func main() {
    // 创建日志器
    logger := grolog.New(nil,
        grolog.WithLevel(grolog.LevelTrace),
        grolog.WithStyle(grolog.StyleDetail),
        grolog.WithEnableAsyn(true),
        grolog.WithFileDir("logs"),
        grolog.WithMaxFileSize(10*grolog.MiB),
    )
	defer logger.Close()

    // 记录日志
    logger.Trace("This is a trace message.\n")
    logger.Warningln("This is a warning message.")
}
```

### JSON 配置支持

Grolog 支持从 JSON 配置文件加载日志配置。

首先, 创建一个 JSON 配置文件,例如 `config.json`, 然后在代码中加载配置:

```go
import (
    "os"
    "grolog"
)

func main() {
    // 日志器配置
	config := grolog.DefaultConfig()
    if js, err := os.ReadFile("config.json"); err != nil {
		fmt.Println("read json file failed,", err)
	} else if err := json.Unmarshal(js, &config); err != nil {
		fmt.Println("read json file failed,", err)
	}
    // 创建日志器
	logger = grolog.New(config)
	defer logger.Close()

	// 将日志配置写入到Json文件
	if js, err := json.MarshalIndent(&config, "", "  "); err != nil {
		fmt.Println("convert options failed", err)
	} else if err = os.WriteFile("config.json", js, 0644); err != nil {
		fmt.Println("write json file failed,", err)
	}

    // 记录日志
    logger.Debug("This is a debug message.\n")
    logger.Warningln("This is a warning message.")
}
```

### 高级用法

Grolog 提供了一些高级用法,例如自定义异常日志处理、日志消息回调和异步执行函数。

#### 异常日志处理

您可以自定义异常日志处理函数,以便在发生致命错误时执行特定操作。使用 `WithFatalHandling` 配置选项进行设置:

```go
import (
    "grolog"
)

func handleFatal(logger *grolog.Logger, r any) {
    // 记录致命错误
    logger.Errorf("Fatal error occurred: %+v", r)

    // 执行其他操作,如写入错误文件、发送警报等
    // ...

    // 关闭日志器
    logger.Close()
    // 终止程序运行
    os.Exit(1)
}

func main() {
	// 创建日志器
    logger := grolog.New(nil,
        grolog.WithFatalHandling(handleFatal),
    )
	defer logger.Close()

    // 记录日志
    logger.Fatal("A fatal error occurred.\n")
}
```

#### 日志消息回调

您可以设置一个回调函数,在记录日志时执行自定义操作。使用 `WithMsgCallback` 配置选项进行设置:

```go
import (
    "grolog"
)

func logMessageCallback(level int, message string) {
    // 执行自定义操作,如发送日志到远程服务器、写入其他文件等
    // ...
}

func main() {
	// 创建日志器
    logger := grolog.New(nil,
        grolog.WithMsgCallback(logMessageCallback),
    )
	defer logger.Close()

    // 记录日志
    logger.Debug("This is a debug message.")
    logger.Warning("This is a warning message.")
}
```

#### 异步执行函数

您可以自定义异步执行函数,以便在执行异步操作时使用。使用 `WithGoExec` 配置选项进行设置:

```go
import (
    "runtime"
    "grolog"
)

func customGoExec(f func()) {
    // 执行自定义操作,如设置 goroutine 属性、限制并发数等
    // ...

    go f()
}

func main() {
    // 设置最大并发数
    runtime.GOMAXPROCS(4)
	// 创建日志器
    logger := grolog.New(nil,
        grolog.WithGoExec(customGoExec),
    )
	defer logger.Close()

    // 记录日志
    logger.Debug("This is a debug message.")
    logger.Warning("This is a warning message.")
}
```
