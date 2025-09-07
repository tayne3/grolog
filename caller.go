// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

// 调用器
type Caller struct {
	handler groHandler // 日志处理器
	layer   int        // 调用层级
}

type groCaller = Caller

func (c groCaller) IsValid() bool {
	return c.handler != nil
}

func (c groCaller) VerBose(a ...any) {
	c.handler.Log(LevelVerBose, c.layer, a...)
}

func (c groCaller) Debug(a ...any) {
	c.handler.Log(LevelDebug, c.layer, a...)
}

func (c groCaller) Trace(a ...any) {
	c.handler.Log(LevelTrace, c.layer, a...)
}

func (c groCaller) Warning(a ...any) {
	c.handler.Log(LevelWarning, c.layer, a...)
}

func (c groCaller) Error(a ...any) {
	c.handler.Log(LevelError, c.layer, a...)
}

func (c groCaller) Fatal(a ...any) {
	c.handler.Log(LevelFatal, c.layer, a...)
}

func (c groCaller) VerBoseln(a ...any) {
	c.handler.Logln(LevelVerBose, c.layer, a...)
}

func (c groCaller) Debugln(a ...any) {
	c.handler.Logln(LevelDebug, c.layer, a...)
}

func (c groCaller) Traceln(a ...any) {
	c.handler.Logln(LevelTrace, c.layer, a...)
}

func (c groCaller) Warningln(a ...any) {
	c.handler.Logln(LevelWarning, c.layer, a...)
}

func (c groCaller) Errorln(a ...any) {
	c.handler.Logln(LevelError, c.layer, a...)
}

func (c groCaller) Fatalln(a ...any) {
	c.handler.Logln(LevelFatal, c.layer, a...)
}

func (c groCaller) VerBosef(format string, args ...any) {
	c.handler.Logf(LevelVerBose, c.layer, format, args...)
}

func (c groCaller) Debugf(format string, args ...any) {
	c.handler.Logf(LevelDebug, c.layer, format, args...)
}

func (c groCaller) Tracef(format string, args ...any) {
	c.handler.Logf(LevelTrace, c.layer, format, args...)
}

func (c groCaller) Warningf(format string, args ...any) {
	c.handler.Logf(LevelWarning, c.layer, format, args...)
}

func (c groCaller) Errorf(format string, args ...any) {
	c.handler.Logf(LevelError, c.layer, format, args...)
}

func (c groCaller) Fatalf(format string, args ...any) {
	c.handler.Logf(LevelFatal, c.layer, format, args...)
}
