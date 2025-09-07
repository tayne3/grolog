// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

// 日志处理器
type groHandler interface {
	Flush()
	Close()
	Log(level int, layer int, a ...any)
	Logln(level int, layer int, a ...any)
	Logf(level int, layer int, format string, args ...any)
}

// 日志器
type Logger struct {
	config  Config     // 日志配置
	handler groHandler // 日志处理器
}

// 创建日志器
func New(cfg *Config, opts ...Option) (l *Logger) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	l = &Logger{config: *cfg}
	for _, opt := range opts {
		opt(&l.config)
	}
	l.config.init(l)

	if l.config.EnableAsyn {
		l.handler = newHandlerAsyn(&l.config)
	} else {
		l.handler = newHandlerSync(&l.config)
	}
	return l
}

// 创建日志器 (使用默认配置)
func Default() (l *Logger) {
	l = &Logger{
		config: *DefaultConfig(),
	}
	l.config.init(l)

	if l.config.EnableAsyn {
		l.handler = newHandlerAsyn(&l.config)
	} else {
		l.handler = newHandlerSync(&l.config)
	}
	return l
}

// 关闭日志器
func (l *Logger) Close() {
	if l.handler != nil {
		l.handler.Close()
	}
}

// 刷新日志器
func (l *Logger) Flush() {
	l.handler.Flush()
}

// 获取调用信息
func (l *Logger) Caller(layer int) Caller {
	return Caller{
		handler: l.handler,
		layer:   layer,
	}
}

func (l *Logger) VerBose(a ...any) {
	l.handler.Log(LevelVerBose, 0, a...)
}

func (l *Logger) Debug(a ...any) {
	l.handler.Log(LevelDebug, 0, a...)
}

func (l *Logger) Trace(a ...any) {
	l.handler.Log(LevelTrace, 0, a...)
}

func (l *Logger) Warning(a ...any) {
	l.handler.Log(LevelWarning, 0, a...)
}

func (l *Logger) Error(a ...any) {
	l.handler.Log(LevelError, 0, a...)
}

func (l *Logger) Fatal(a ...any) {
	l.handler.Log(LevelFatal, 0, a...)
}

func (l *Logger) VerBoseln(a ...any) {
	l.handler.Logln(LevelVerBose, 0, a...)
}

func (l *Logger) Debugln(a ...any) {
	l.handler.Logln(LevelDebug, 0, a...)
}

func (l *Logger) Traceln(a ...any) {
	l.handler.Logln(LevelTrace, 0, a...)
}

func (l *Logger) Warningln(a ...any) {
	l.handler.Logln(LevelWarning, 0, a...)
}

func (l *Logger) Errorln(a ...any) {
	l.handler.Logln(LevelError, 0, a...)
}

func (l *Logger) Fatalln(a ...any) {
	l.handler.Logln(LevelFatal, 0, a...)
}

func (l *Logger) VerBosef(format string, args ...any) {
	l.handler.Logf(LevelVerBose, 0, format, args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.handler.Logf(LevelDebug, 0, format, args...)
}

func (l *Logger) Tracef(format string, args ...any) {
	l.handler.Logf(LevelTrace, 0, format, args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.handler.Logf(LevelWarning, 0, format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.handler.Logf(LevelError, 0, format, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.handler.Logf(LevelFatal, 0, format, args...)
}
