// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

const (
	LevelVerBose int = iota // 变量信息, 用于记录变量信息
	LevelDebug              // 调试信息, 用于打印详细的调试信息
	LevelTrace              // 提示信息, 用于跟踪程序执行流程
	LevelWarning            // 警告信息, 用于打印警告信息
	LevelError              // 错误信息, 用于记录非致命性的错误信息
	LevelFatal              // 致命错误, 用于记录致命的错误信息
)

// 日志级别字符串
var levelStrings = []string{
	LevelVerBose: "VBOSE",
	LevelDebug:   "DEBUG",
	LevelTrace:   "TRACE",
	LevelWarning: "WARNG",
	LevelError:   "ERROR",
	LevelFatal:   "FATAL",
}

// 日志样式字符串-起始
var levelStyleStarts = []string{
	LevelVerBose: "\x1b[32;2m",
	LevelDebug:   "\x1b[32;2m",
	LevelTrace:   "\x1b[36;2m",
	LevelWarning: "\x1b[33;2m",
	LevelError:   "\x1b[31;2m",
	LevelFatal:   "\x1b[31;2m",
}

// 日志样式字符串-结束
const (
	levelStyleEnd = "\x1b[0m"
)

const (
	StyleBasic  int = iota // 基本日志 (只包含消息)
	StyleBrief             // 简要日志 (默认日志样式, 包含 级别、时间、消息)
	StyleDetail            // 调试日志 (会影响性能, 包含 级别、时间、文件名、行号、消息)
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
)

const (
	defaultLevel          = LevelWarning // 默认日志级别
	defaultStyle          = StyleBrief   // 默认日志样式
	defaultMaxAsynExec    = 100          // 默认异步执行数量上限
	defaultMaxAsynBuffer  = 128          // 默认异步缓冲大小
	defaultMaxWriteBuffer = 4096         // 默认日志文件缓冲大小
	defaultMaxFileCount   = 5            // 默认日志文件数量上限
	defaultMaxFileSize    = 5 * MiB      // 默认日志文件大小上限
	defaultFileDir        = "log"        // 默认日志文件保存目录
	defaultFlashInterval  = "3h0m0s"     // 默认日志文件刷新间隔 (3h)
	defaultExpireTime     = "0s"         // 默认日志文件过期时间 (默认禁用)
)

// 定义配置选项
type Config struct {
	logger         *Logger            `json:"-"`              // 日志器 (永不为空)
	startTime      time.Time          `json:"-"`              // 启始时间 (创建时自动填充)
	FatalHandling  func(*Logger, any) `json:"-"`              // 异常日志处理函数 (为空时使用默认值)
	MsgCallback    func(int, string)  `json:"-"`              // 日志消息回调函数 (默认为空)
	GoExec         func(func())       `json:"-"`              // 异步执行函数 (为空时使用go语句执行)
	Level          int                `json:"Level"`          // 日志级别 (默认警告级别, 值无效时使用默认值)
	Style          int                `json:"Style"`          // 日志样式 (默认简要样式, 值无效时使用默认值)
	EnableAsyn     bool               `json:"EnableAsyn"`     // 是否启用异步模式 (默认禁用异步模式, 同步模式的性能可能会优于异步模式，但异步模式下资源使用更加可控)
	EnableFileTime bool               `json:"EnableFileTime"` // 是否启用文件时间 (默认禁用文件名包含时间信息)
	DisableSave    bool               `json:"DisableSave"`    // 是否禁用日志文件 (默认启用日志文件)
	DisablePrint   bool               `json:"DisablePrint"`   // 是否禁用日志打印 (默认启用日志打印)
	MaxAsynExec    int                `json:"MaxAsynExec"`    // 异步执行数量上限 (启用异步模式时有效, 小于等于0时使用默认值)
	MaxAsynBuffer  int                `json:"MaxAsynBuffer"`  // 异步消息缓冲大小 (启用异步模式时有效, 小于0时使用默认值)
	MaxWriteBuffer int                `json:"MaxWriteBuffer"` // 日志文件缓冲大小 (启用日志文件时有效, 小于0时使用默认值)
	MaxFileCount   int                `json:"MaxFileCount"`   // 日志文件数量上限 (启用日志文件时有效, 小于等于0时使用默认值)
	MaxFileSize    int64              `json:"MaxFileSize"`    // 日志文件大小上限 (启用日志文件时有效, 小于等于0时使用默认值)
	FileDir        string             `json:"FileDir"`        // 日志文件保存目录 (启用日志文件时有效, 为空时使用程序运行路径下的log目录)
	FileName       string             `json:"FileName"`       // 日志文件保存名称 (启用日志文件时有效, 为空时使用程序名称)
	FlashInterval  string             `json:"FlashInterval"`  // 日志文件刷新间隔 (启用日志文件时有效, 等于0时禁用, 值无效时使用默认值)
	ExpireTime     string             `json:"ExpireTime"`     // 日志文件过期时间 (启用日志文件时有效, 等于0时禁用, 值无效时使用默认值)
}

// 配置选项
type Option func(*Config)

// 默认配置
func DefaultConfig() *Config {
	return &Config{
		logger:         nil,
		startTime:      time.Now(),
		FatalHandling:  nil,
		MsgCallback:    nil,
		GoExec:         nil,
		Level:          defaultLevel,
		Style:          defaultStyle,
		EnableAsyn:     false,
		EnableFileTime: false,
		DisableSave:    false,
		DisablePrint:   false,
		MaxAsynExec:    defaultMaxAsynExec,
		MaxAsynBuffer:  defaultMaxAsynBuffer,
		MaxWriteBuffer: defaultMaxWriteBuffer,
		MaxFileCount:   defaultMaxFileCount,
		MaxFileSize:    defaultMaxFileSize,
		FileDir:        defaultFileDir,
		FileName:       "",
		FlashInterval:  defaultFlashInterval,
		ExpireTime:     defaultExpireTime,
	}
}

// 初始化配置
func (c *Config) init(logger *Logger) {
	c.logger = logger
	if c.FatalHandling == nil {
		c.FatalHandling = func(l *Logger, r any) {
			l.Error("abnormal exits from panic: ", r, "\n", string(debug.Stack()))
			l.Close()
			os.Exit(1)
		}
	}
	if c.GoExec == nil {
		c.GoExec = func(f func()) {
			go f()
		}
	}
	if c.Level < LevelVerBose || c.Level > LevelFatal {
		c.Level = defaultLevel
	}
	if c.Style < StyleBasic || c.Style > StyleDetail {
		c.Style = defaultStyle
	}
	if c.MaxAsynExec <= 0 {
		c.MaxAsynExec = defaultMaxAsynExec
	}
	if c.MaxAsynBuffer < 0 {
		c.MaxAsynBuffer = defaultMaxAsynBuffer
	}
	if c.MaxWriteBuffer < 0 {
		c.MaxWriteBuffer = defaultMaxWriteBuffer
	}
	if c.MaxFileCount <= 0 {
		c.MaxFileCount = defaultMaxFileCount
	}
	if c.MaxFileSize <= 0 {
		c.MaxFileSize = defaultMaxFileSize
	}
	if c.FileDir == "" {
		path, _ := os.Executable()
		c.FileDir = filepath.Join(filepath.Dir(path), "log")
	}
	if c.FileName == "" {
		if path, err := os.Executable(); err != nil {
			panic(err)
		} else {
			path = filepath.Base(path)
			c.FileName = strings.TrimSuffix(path, filepath.Ext(path))
		}
	}
	if duration, err := time.ParseDuration(c.FlashInterval); err != nil || duration < 0 {
		c.FlashInterval = defaultFlashInterval
	}
	if duration, err := time.ParseDuration(c.ExpireTime); err != nil || duration < 0 {
		c.ExpireTime = defaultExpireTime
	}
}

// 使用配置选项
func (c *Config) Use(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

// 设置异常日志处理
func WithFatalHandling(handling func(*Logger, any)) Option {
	return func(opt *Config) {
		opt.FatalHandling = handling
	}
}

// 设置日志消息回调
func WithMsgCallback(handler func(int, string)) Option {
	return func(opt *Config) {
		opt.MsgCallback = handler
	}
}

// 设置异步执行函数
func WithGoExec(exec func(f func())) Option {
	return func(opt *Config) {
		opt.GoExec = exec
	}
}

// 设置日志级别
func WithLevel(level int) Option {
	return func(opt *Config) {
		opt.Level = level
	}
}

// 设置日志样式
func WithStyle(style int) Option {
	return func(opt *Config) {
		opt.Style = style
	}
}

// 设置是否启用异步
func WithEnableAsyn(asyn bool) Option {
	return func(opt *Config) {
		opt.EnableAsyn = asyn
	}
}

// 设置是否启用文件时间
func WithEnableFileTime(enable bool) Option {
	return func(opt *Config) {
		opt.EnableFileTime = enable
	}
}

// 设置是否禁用日志文件
func WithDisableSave(save bool) Option {
	return func(opt *Config) {
		opt.DisableSave = save
	}
}

// 设置是否禁用日志打印
func WithDisablePrint(print bool) Option {
	return func(opt *Config) {
		opt.DisablePrint = print
	}
}

// 设置异步执行数量上限
func WithAsynMaxGor(max int) Option {
	return func(opt *Config) {
		opt.MaxAsynExec = max
	}
}

// 设置异步消息缓冲大小
func WithAsynMaxBuffer(max int) Option {
	return func(opt *Config) {
		opt.MaxAsynBuffer = max
	}
}

// 设置日志文件缓冲大小
func WithWriteBufferSize(size int) Option {
	return func(opt *Config) {
		opt.MaxWriteBuffer = size
	}
}

// 设置日志文件保存目录
func WithFileDir(dir string) Option {
	return func(opt *Config) {
		opt.FileDir = dir
	}
}

// 设置日志文件保存名称
func WithFileName(name string) Option {
	return func(opt *Config) {
		opt.FileName = name
	}
}

// 设置日志文件大小上限
func WithMaxFileSize(maxSize int64) Option {
	return func(opt *Config) {
		opt.MaxFileSize = maxSize
	}
}

// 设置日志文件数量上限
func WithMaxFileCount(maxCount int) Option {
	return func(opt *Config) {
		opt.MaxFileCount = maxCount
	}
}

// 设置日志文件刷新间隔
func WithFlashInterval(interval string) Option {
	return func(opt *Config) {
		opt.FlashInterval = interval
	}
}

// 设置日志文件过期时间
func WithExpireTime(expire string) Option {
	return func(opt *Config) {
		opt.ExpireTime = expire
	}
}
