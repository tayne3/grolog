// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 日志处理器-异步
type groHandlerAsyn struct {
	config      *Config            // 日志选项 (永不为空)
	pusher      *groPusher         // 日志推送器 (永不为空)
	flashCancel context.CancelFunc // 定时刷新结束
	execCount   atomic.Int32       // 异步执行计数
	closed      atomic.Bool        // 是否已关闭
	msgs        chan *groMsg       // 消息队列
	stop        sync.WaitGroup     // 停止等待
}

var _ groHandler = (*groHandlerAsyn)(nil)

// 创建异步日志处理器
func newHandlerAsyn(config *Config) groHandler {
	h := &groHandlerAsyn{
		config:    config,
		pusher:    newPusher(config),
		execCount: atomic.Int32{},
		closed:    atomic.Bool{},
		msgs:      make(chan *groMsg, config.MaxAsynBuffer),
		stop:      sync.WaitGroup{},
	}

	// 预分配消息对象
	if config.MaxAsynBuffer > 0 {
		msgs := make([]*groMsg, config.MaxAsynBuffer)
		for i := 0; i < config.MaxAsynBuffer; i++ {
			msgs[i] = h.pusher.get()
		}
		for i := 0; i < config.MaxAsynBuffer; i++ {
			h.pusher.put(msgs[i])
		}
	}

	// 定时刷新缓冲区
	duration, _ := time.ParseDuration(h.config.FlashInterval)
	var ctx context.Context
	ctx, h.flashCancel = context.WithCancel(context.Background())
	h.goFlash(duration, ctx)

	return h
}

// 关闭日志处理器
func (h *groHandlerAsyn) Close() {
	if !h.closed.CompareAndSwap(false, true) {
		h.stop.Wait()
		h.pusher.Close()
		return
	}

	h.flashCancel()
	h.msgs <- nil // 未关闭的话, 发送空消息
	h.stop.Wait()
	h.pusher.Close()
	close(h.msgs)
}

// 刷新日志处理器
func (h *groHandlerAsyn) Flush() {
	if h.closed.Load() {
		return
	}

	h.msgs <- nil
}

func (h *groHandlerAsyn) Log(level int, layer int, a ...any) {
	if h.config.Level > level || h.closed.Load() {
		return
	}

	m := h.pusher.get()
	h.pusher.assign(m, level, layer)
	fmt.Fprint(m.text, a...)

	h.msgHanding(m)
}

func (h *groHandlerAsyn) Logln(level int, layer int, a ...any) {
	if h.config.Level > level || h.closed.Load() {
		return
	}

	m := h.pusher.get()
	h.pusher.assign(m, level, layer)
	fmt.Fprintln(m.text, a...)

	h.msgHanding(m)
}

func (h *groHandlerAsyn) Logf(level int, layer int, format string, args ...any) {
	if h.config.Level > level || h.closed.Load() {
		return
	}

	m := h.pusher.get()
	h.pusher.assign(m, level, layer)
	fmt.Fprintf(m.text, format, args...)

	h.msgHanding(m)
}

// 消息处理
func (h *groHandlerAsyn) msgHanding(m *groMsg) {
	retry := 0
	for {
		select {
		case h.msgs <- m:
			if h.execCount.Load() == 0 {
				h.execCount.Add(1)
				h.goHanding()
			}
			return
		default:
			if h.closed.Load() {
				return
			}
			retry++
			if retry > 3 && h.execCount.Load() < int32(h.config.MaxAsynExec) {
				h.execCount.Add(1)
				h.goHanding()
			}
			runtime.Gosched()
		}
	}
}

// 消息处理
func (h *groHandlerAsyn) goHanding() {
	h.stop.Add(1)
	go func() {
		defer h.stop.Done()

		idle := 0
		retry := 0
		for {
			select {
			case m := <-h.msgs: // 不断取出消息
				idle = 0
				retry++
				if retry > 100 {
					runtime.Gosched()
				}
				if m != nil {
					h.pusher.push(m)
					h.pusher.put(m)
					break
				}
				if h.closed.Load() { // 收到空消息且已关闭, 结束
					goto Closed
				}
				h.pusher.Flush() // 未关闭, 刷新缓冲区
			default:
				retry = 0
				idle++
				if idle > 100 && h.execCount.Load() > 0 {
					h.execCount.Add(-1)
					goto Closed
				}
				// runtime.Gosched()
			}
		}
	Closed:
		if !h.closed.Load() {
			return
		}
		for {
			select {
			case m := <-h.msgs: // 取出已有的消息
				if m != nil {
					h.pusher.push(m)
					h.pusher.put(m)
				}
			default: // 无消息, 结束
				h.msgs <- nil
				return
			}
		}
	}()
}

// 定时刷新日志
func (h *groHandlerAsyn) goFlash(interval time.Duration, ctx context.Context) {
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
