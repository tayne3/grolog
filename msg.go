// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 日志消息
type groMsg struct {
	level int
	tips  *bytes.Buffer
	stack *bytes.Buffer
	text  *bytes.Buffer
}

// 填充基本日志消息
func (m *groMsg) initBasic(int) {
}

// 填充简要日志消息
func (m *groMsg) initBrief(int) {
	m.tips.WriteString(levelStrings[m.level])
	m.tips.WriteString("|")
	{
		time := time.Now().Format("06.01.02-15:04:05.000")
		m.tips.WriteString(time)
	}
}

// 填充详细日志消息
func (m *groMsg) initDetailed(layer int) {
	m.tips.WriteString(levelStrings[m.level])
	m.tips.WriteString("|")
	{
		time := time.Now().Format("06.01.02-15:04:05.000")
		m.tips.WriteString(time)
	}

	m.stack.WriteString("[")
	if _, file, line, ok := runtime.Caller(4 + layer); !ok {
		m.stack.WriteString("???#?")
	} else {
		m.stack.WriteString(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)))
		m.stack.WriteString(":")
		m.stack.WriteString(strconv.Itoa(line))
	}
	m.stack.WriteString("]")
}

// 缓存对象池
type groBufferPool struct {
	pool sync.Pool
}

// 初始化对象池
func (p *groBufferPool) Init() {
	p.pool.New = func() any {
		return new(bytes.Buffer)
	}
}

// 从对象池获取
func (p *groBufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// 回收对象
func (p *groBufferPool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}

// 消息对象池
type groMsgPool struct {
	pool sync.Pool
}

// 初始化对象池
func (p *groMsgPool) Init() {
	p.pool.New = func() any {
		return new(groMsg)
	}
}

// 从对象池获取
func (p *groMsgPool) Get() *groMsg {
	return p.pool.Get().(*groMsg)
}

// 回收对象
func (p *groMsgPool) Put(m *groMsg) {
	p.pool.Put(m)
}
