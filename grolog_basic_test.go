// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	testTask   = 10
	testNumber = 1000
)

var (
	curMem     = uint64(0)
	testCount  int
	lock       sync.Mutex
	msgHanlder = func(int, string) {
		lock.Lock()
		testCount++
		lock.Unlock()
		time.Sleep(5 * time.Millisecond)
	}
)

func TestSyncOutBasic(t *testing.T) {
	config := []Option{
		WithLevel(LevelVerBose),
		WithStyle(StyleBasic),
		WithEnableAsyn(false),
		WithDisableSave(true),
		WithDisablePrint(false),
		WithMsgCallback(msgHanlder),
	}

	testOut(t, "Sync-Basic", testTask, testNumber, config...)
}

func TestSyncOutBrief(t *testing.T) {
	config := []Option{
		WithLevel(LevelVerBose),
		WithStyle(StyleBrief),
		WithEnableAsyn(false),
		WithDisableSave(true),
		WithDisablePrint(false),
		WithMsgCallback(msgHanlder),
	}

	testOut(t, "Sync-Brief", testTask, testNumber, config...)
}

func TestSyncOutDetail(t *testing.T) {
	config := []Option{
		WithLevel(LevelVerBose),
		WithStyle(StyleDetail),
		WithEnableAsyn(false),
		WithDisableSave(true),
		WithDisablePrint(false),
		WithMsgCallback(msgHanlder),
	}

	testOut(t, "Sync-Detail", testTask, testNumber, config...)
}

func TestAsynOutBasic(t *testing.T) {
	config := []Option{
		WithLevel(LevelVerBose),
		WithStyle(StyleBasic),
		WithEnableAsyn(true),
		WithDisableSave(true),
		WithDisablePrint(false),
		WithAsynMaxBuffer(128),
		WithMsgCallback(msgHanlder),
	}

	testOut(t, "Asyn-Basic", testTask, testNumber, config...)
}

func TestAsynOutBrief(t *testing.T) {
	config := []Option{
		WithLevel(LevelVerBose),
		WithStyle(StyleBrief),
		WithEnableAsyn(true),
		WithDisableSave(true),
		WithDisablePrint(false),
		WithAsynMaxBuffer(128),
		WithMsgCallback(msgHanlder),
	}
	testOut(t, "Asyn-Brief", testTask, testNumber, config...)
}

func TestAsynOutDetail(t *testing.T) {
	config := []Option{
		WithLevel(LevelVerBose),
		WithStyle(StyleDetail),
		WithEnableAsyn(true),
		WithDisableSave(true),
		WithDisablePrint(false),
		WithAsynMaxBuffer(128),
		WithMsgCallback(msgHanlder),
	}
	testOut(t, "Asyn-Detail", testTask, testNumber, config...)
}

func testOut(t *testing.T, title string, task int, number int, config ...Option) {
	testCount = 0
	logger := New(nil, config...)
	wg := sync.WaitGroup{}
	startTime := time.Now()

	for t := 0; t < task; t++ {
		wg.Add(1)
		go func() {
			for n := 0; n < number; n++ {
				logger.Trace("Hello World!\n")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	submitTime := time.Now()
	logger.Close()
	execTime := time.Now()

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/KiB - curMem

	fmt.Printf("[%s] \n", title)
	fmt.Printf(" - task number   : %d \n", task*number)
	fmt.Printf(" - test count    : %d \n", testCount)
	fmt.Printf(" - submit spent  : %s \n", submitTime.Sub(startTime))
	fmt.Printf(" - execute spent : %s \n", execTime.Sub(startTime))
	fmt.Printf(" - memory usage  : %d KB \n", curMem)

	if testCount != task*number {
		t.Errorf("Test failed: %s", title)
	}
}
