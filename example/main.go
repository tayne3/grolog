// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"grolog"
	"os"
	"runtime/debug"
	"time"
)

var logger *grolog.Logger

func ExDebug(a ...any) {
	logger.Caller(1).Debug(a...)
}

func main() {
	// 创建日志器
	logger = grolog.New(nil,
		grolog.WithLevel(grolog.LevelVerBose),
		grolog.WithStyle(grolog.StyleBrief),
		grolog.WithMsgCallback(nil),
		grolog.WithGoExec(func(f func()) {
			go f()
		}),
		grolog.WithFatalHandling(func(l *grolog.Logger, r any) { // 重定向异常处理
			l.Error("abnormal exits from panic: ", r, "\n", string(debug.Stack()))
			l.Close()
			os.Exit(1)
		}),
	)
	defer logger.Close()

	for i := 0; i < 6; i++ {
		switch i {
		case 0:
			logger.VerBose("VerBose - ", i, ", now time: ", time.Now().Format("2006-01-02 15:04:05"), "\n")
		case 1:
			ExDebug("Debug - ", i, ", now time: ", time.Now().Format("2006-01-02 15:04:05"), "\n")
		case 2:
			logger.Traceln("Traceln -", i, ", now time: ", time.Now().Format("2006-01-02 15:04:05"))
		case 3:
			logger.Warningln("Warningln -", i, ", now time: ", time.Now().Format("2006-01-02 15:04:05"))
		case 4:
			logger.Errorf("Errorf - %d, now time: %s\n", i, time.Now().Format("2006-01-02 15:04:05"))
		case 5:
			logger.Fatal("Fatal - ", i, ", now time: ", time.Now().Format("2006-01-02 15:04:05"), "\n")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
