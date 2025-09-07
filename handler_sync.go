// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"
)

// 日志处理器-同步
type groHandlerSync struct {
	config      *Config            // 日志选项 (永不为空)
	pusher      *groPusher         // 日志推送器 (永不为空)
	stop        sync.WaitGroup     // 停止等待
	flashCancel context.CancelFunc // 定时刷新结束
	bufferPool  sync.Pool          // 缓冲区对象池
}

var _ groHandler = (*groHandlerSync)(nil)

// 创建同步日志处理器
func newHandlerSync(config *Config) groHandler {
	h := &groHandlerSync{
		config:     config,
		stop:       sync.WaitGroup{},
		pusher:     newPusher(config),
		bufferPool: sync.Pool{New: func() any { return new(bytes.Buffer) }},
	}

	duration, _ := time.ParseDuration(h.config.FlashInterval)
	var ctx context.Context
	ctx, h.flashCancel = context.WithCancel(context.Background())
	h.goFlash(duration, ctx)

	return h
}

// 关闭日志处理器
func (h *groHandlerSync) Close() {
	h.flashCancel()
	h.stop.Wait()
	h.pusher.Close()
}

// 刷新日志处理器
func (h *groHandlerSync) Flush() {
	h.pusher.Flush()
}

func (h *groHandlerSync) Log(level int, layer int, a ...any) {
	if h.config.Level > level {
		return
	}

	var m groMsg
	h.pusher.assign(&m, level, layer)
	fmt.Fprint(m.text, a...)

	h.pusher.push(&m)
}

func (h *groHandlerSync) Logln(level int, layer int, a ...any) {
	if h.config.Level > level {
		return
	}

	var m groMsg
	h.pusher.assign(&m, level, layer)
	fmt.Fprintln(m.text, a...)

	h.pusher.push(&m)
}

func (h *groHandlerSync) Logf(level int, layer int, format string, args ...any) {
	if h.config.Level > level {
		return
	}

	var m groMsg
	h.pusher.assign(&m, level, layer)
	fmt.Fprintf(m.text, format, args...)

	h.pusher.push(&m)
}

// 定时刷新日志
func (h *groHandlerSync) goFlash(interval time.Duration, ctx context.Context) {
	if interval <= 0 {
		return
	}

	h.stop.Add(1)
	go func() {
		defer h.stop.Done()
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			h.Flush()
		}
	}()
}
