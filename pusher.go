// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"os"
	"sync"
)

// 日志推送器
type groPusher struct {
	config     *Config       // 日志选项
	out        *os.File      // 打印输出
	storage    *groStorage   // 日志存储
	closed     bool          // 是否已关闭
	msgPool    groMsgPool    // 消息对象池
	bufferPool groBufferPool // 缓冲区对象池
	lock       sync.Mutex    // 打印锁
}

// 创建新的推送器
func newPusher(config *Config) *groPusher {
	p := &groPusher{
		config:     config,
		out:        nil,
		storage:    nil,
		closed:     false,
		msgPool:    groMsgPool{},
		bufferPool: groBufferPool{},
	}

	if !p.config.DisablePrint {
		p.out = os.Stdout
	}
	if !p.config.DisableSave {
		p.storage = newStorage(config)
	}

	p.msgPool.Init()
	p.bufferPool.Init()

	return p
}

// 刷新推送器
func (p *groPusher) Flush() {
	if p.closed {
		return
	}
	if p.out != nil {
		p.out.Sync()
	}
	if p.storage != nil {
		p.storage.Flush()
	}
}

// 关闭推送器
func (p *groPusher) Close() {
	if p.closed {
		return
	}
	if p.storage != nil {
		p.storage.Close()
	}
	p.closed = true
}

// 获取消息对象
func (p *groPusher) get() *groMsg {
	return p.msgPool.Get()
}

// 回收消息对象
func (p *groPusher) put(m *groMsg) {
	p.msgPool.Put(m)
}

// 填充消息
func (p *groPusher) assign(m *groMsg, level int, layer int) {
	if m.level != level {
		m.level = level
	}

	switch p.config.Style {
	case StyleBasic:
		m.text = p.bufferPool.Get()
		m.text.Reset()
		m.initBasic(layer)
	case StyleBrief:
		m.text = p.bufferPool.Get()
		m.tips = p.bufferPool.Get()
		m.text.Reset()
		m.tips.Reset()
		m.initBrief(layer)
	case StyleDetail:
		m.text = p.bufferPool.Get()
		m.tips = p.bufferPool.Get()
		m.stack = p.bufferPool.Get()
		m.text.Reset()
		m.tips.Reset()
		m.stack.Reset()
		m.initDetailed(layer)
	}
}

// 推送日志消息
func (p *groPusher) push(m *groMsg) {
	if p.closed {
		return
	}

	if p.out != nil {
		switch p.config.Style {
		case StyleBasic:
			p.out.Write(m.text.Bytes())
		case StyleBrief:
			p.lock.Lock()
			p.out.WriteString(levelStyleStarts[m.level])
			p.out.Write(m.tips.Bytes())
			p.out.WriteString(levelStyleEnd)
			p.out.WriteString(" ")
			p.out.Write(m.text.Bytes())
			p.lock.Unlock()
		case StyleDetail:
			p.lock.Lock()
			p.out.WriteString(levelStyleStarts[m.level])
			p.out.Write(m.tips.Bytes())
			p.out.WriteString(levelStyleEnd)
			p.out.WriteString(" ")
			p.out.Write(m.stack.Bytes())
			p.out.WriteString(" ")
			p.out.Write(m.text.Bytes())
			p.lock.Unlock()
		}
	}

	if p.storage != nil {
		switch p.config.Style {
		case StyleBasic:
			p.storage.Write(m.text.Bytes())
		case StyleBrief:
			p.lock.Lock()
			p.storage.Write(m.tips.Bytes())
			p.storage.WriteString(" ")
			p.storage.Write(m.text.Bytes())
			p.lock.Unlock()
		case StyleDetail:
			p.lock.Lock()
			p.storage.Write(m.tips.Bytes())
			p.storage.WriteString(" ")
			p.storage.Write(m.text.Bytes())
			p.lock.Unlock()
		}
	}

	if p.config.MsgCallback != nil {
		p.config.GoExec(func() {
			p.config.MsgCallback(m.level, m.text.String())
		})
	}

	switch p.config.Style {
	case StyleBasic:
		m.text.Reset()
		p.bufferPool.Put(m.text)
	case StyleBrief:
		m.tips.Reset()
		m.text.Reset()
		p.bufferPool.Put(m.tips)
		p.bufferPool.Put(m.text)
	case StyleDetail:
		m.tips.Reset()
		m.text.Reset()
		m.stack.Reset()
		p.bufferPool.Put(m.tips)
		p.bufferPool.Put(m.text)
		p.bufferPool.Put(m.stack)
	}

	if m.level == LevelFatal {
		go p.config.FatalHandling(p.config.logger, recover())
	}
}
