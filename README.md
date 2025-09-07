<p align="center" style="font-size: 24px; font-weight: bold; color: #797979;">Grolog</p>

**English** | [中文](README_zh.md)

Grolog is a lightweight Go logging library designed to provide simple, high-performance, and configurable logging functionality. It supports multiple log levels, log formats, and both asynchronous and synchronous modes. In addition, Grolog offers flexible configuration options, allowing you to customize logging behavior to your needs.

## Features

- **Multiple Log Levels**: Verbose, Debug, Trace, Warning, Error, and Fatal.
- **Multiple Log Formats**: Basic, Brief, and Detail.
- **Asynchronous and Synchronous Modes**: Supports both asynchronous and synchronous logging modes, which can be selected as needed.
- **File Logging**: Supports writing logs to files, with configurable file size limits, file count limits, and expiration times.
- **Console Logging**: Supports outputting logs to the console.
- **Highly Configurable**: Provides a variety of configuration options to customize logging behavior as needed.
- **JSON Configuration Support**: Supports loading log configurations from a JSON configuration file.

## Quick Start

### Basic Usage

```go
import (
    "grolog"
)

func main() {
    // Create a logger
    logger := grolog.Default() // Equivalent to `logger := grolog.New(nil)`
	defer logger.Close()

    // Log messages
    logger.Debug("This is a debug message.\n")
    logger.Warningf("Warning: %s\n", "Something went wrong.")
    logger.Errorln("An error occurred.")
}
```

### Configuration Options

Grolog provides a variety of configuration options to customize logging behavior as needed. Here are some common configuration options:

- `WithFatalHandling(handling func(*Logger, any))`: Sets a fatal log handling function to be executed when a fatal error occurs.
- `WithMsgCallback(handler func(int, string))`: Sets a log message callback function to be executed when logging a message.
- `WithGoExec(exec func(f func()))`: Sets an asynchronous execution function to support an external goroutine pool.
- `WithLevel(level int)`: Sets the log level, with possible values of `LevelVerBose`, `LevelDebug`, `LevelTrace`, `LevelWarning`, `LevelError`, and `LevelFatal`.
- `WithStyle(style int)`: Sets the log format, with possible values of `StyleBasic`, `StyleBrief`, and `StyleDetail`.
- `WithEnableAsyn(asyn bool)`: Enables or disables asynchronous logging mode (synchronous mode may have better performance, but asynchronous mode has more controllable resource usage).
- `WithEnableFileTime(enable bool)`: Enables or disables log filenames that include time information.
- `WithDisableSave(save bool)`: Disables or enables file logging.
- `WithDisablePrint(print bool)`: Disables or enables console logging.
- `WithAsynMaxGor(max int)`: Sets the maximum number of asynchronous executions, effective when asynchronous mode is enabled.
- `WithAsynMaxBuffer(max int)`: Sets the asynchronous message buffer size, effective when asynchronous mode is enabled.
- `WithWriteBufferSize(size int)`: Sets the log file buffer size, effective when file logging is enabled.
- `WithFileDir(dir string)`: Sets the log file directory, effective when file logging is enabled.
- `WithFileName(name string)`: Sets the log file name, effective when file logging is enabled.
- `WithMaxFileSize(maxSize int64)`: Sets the maximum size of a single log file, effective when file logging is enabled.
- `WithMaxFileCount(maxCount int)`: Sets the maximum number of log files, effective when file logging is enabled.
- `WithFlashInterval(interval string)`: Sets the log flash interval, effective when file logging is enabled.
- `WithExpireTime(expire string)`: Sets the log file expiration time, effective when file logging is enabled.

Example:

```go
import (
    "grolog"
)

func main() {
    // Create a logger
    logger := grolog.New(nil,
        grolog.WithLevel(grolog.LevelTrace),
        grolog.WithStyle(grolog.StyleDetail),
        grolog.WithEnableAsyn(true),
        grolog.WithFileDir("logs"),
        grolog.WithMaxFileSize(10*grolog.MiB),
    )
	defer logger.Close()

    // Log messages
    logger.Trace("This is a trace message.\n")
    logger.Warningln("This is a warning message.")
}
```

### JSON Configuration Support

Grolog supports loading log configurations from a JSON configuration file.

First, create a JSON configuration file, for example `config.json`, and then load the configuration in your code:

```go
import (
    "os"
    "grolog"
)

func main() {
    // Logger configuration
	config := grolog.DefaultConfig()
    if js, err := os.ReadFile("config.json"); err != nil {
		fmt.Println("read json file failed,", err)
	} else if err := json.Unmarshal(js, &config); err != nil {
		fmt.Println("read json file failed,", err)
	}
    // Create a logger
	logger = grolog.New(config)
	defer logger.Close()

	// Write the log configuration to a JSON file
	if js, err := json.MarshalIndent(&config, "", "  "); err != nil {
		fmt.Println("convert options failed", err)
	} else if err = os.WriteFile("config.json", js, 0644); err != nil {
		fmt.Println("write json file failed,", err)
	}

    // Log messages
    logger.Debug("This is a debug message.\n")
    logger.Warningln("This is a warning message.")
}
```

### Advanced Usage

Grolog provides some advanced usage, such as custom fatal log handling, log message callbacks, and asynchronous execution functions.

#### Fatal Log Handling

You can customize the fatal log handling function to perform specific actions when a fatal error occurs. Use the `WithFatalHandling` configuration option to set it:

```go
import (
    "grolog"
)

func handleFatal(logger *grolog.Logger, r any) {
    // Log the fatal error
    logger.Errorf("Fatal error occurred: %+v", r)

    // Perform other actions, such as writing to an error file, sending alerts, etc.
    // ...

    // Close the logger
    logger.Close()
    // Terminate the program
    os.Exit(1)
}

func main() {
	// Create a logger
    logger := grolog.New(nil,
        grolog.WithFatalHandling(handleFatal),
    )
	defer logger.Close()

    // Log a fatal error
    logger.Fatal("A fatal error occurred.\n")
}
```

#### Log Message Callback

You can set a callback function to perform custom actions when logging a message. Use the `WithMsgCallback` configuration option to set it:

```go
import (
    "grolog"
)

func logMessageCallback(level int, message string) {
    // Perform custom actions, such as sending logs to a remote server, writing to other files, etc.
    // ...
}

func main() {
	// Create a logger
    logger := grolog.New(nil,
        grolog.WithMsgCallback(logMessageCallback),
    )
	defer logger.Close()

    // Log messages
    logger.Debug("This is a debug message.")
    logger.Warning("This is a warning message.")
}
```

#### Asynchronous Execution Function

You can customize the asynchronous execution function to be used when performing asynchronous operations. Use the `WithGoExec` configuration option to set it:

```go
import (
    "runtime"
    "grolog"
)

func customGoExec(f func()) {
    // Perform custom actions, such as setting goroutine properties, limiting concurrency, etc.
    // ...

    go f()
}

func main() {
    // Set the maximum number of concurrent executions
    runtime.GOMAXPROCS(4)
	// Create a logger
    logger := grolog.New(nil,
        grolog.WithGoExec(customGoExec),
    )
	defer logger.Close()

    // Log messages
    logger.Debug("This is a debug message.")
    logger.Warning("This is a warning message.")
}
```
