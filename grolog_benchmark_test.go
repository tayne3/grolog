// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"sync"
	"testing"
)

const (
	RunTimes = 10
)

func BenchmarkLoggerSync(b *testing.B) {
	logger := New(nil,
		WithLevel(LevelVerBose),
		WithStyle(StyleBasic),
		WithDisableSave(true),
		WithDisablePrint(false),
	)
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				logger.Traceln("Hello World!")
				wg.Done()
			}()
		}
		wg.Wait()
	}
	logger.Close()
}

func BenchmarkLoggerAsyn(b *testing.B) {
	logger := New(nil,
		WithLevel(LevelVerBose),
		WithStyle(StyleBasic),
		WithDisableSave(true),
		WithDisablePrint(false),
		WithEnableAsyn(true),
	)
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				logger.Traceln("Hello World!")
				wg.Done()
			}()
		}
		wg.Wait()
	}
	logger.Close()
}
