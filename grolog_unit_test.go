// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestCleanExpiredFiles(t *testing.T) {
	testDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	const (
		testFileCount = 5
	)

	testName := "Test"
	expireTime := 30 * time.Minute
	now := time.Now()
	expired := now.Add(-expireTime)
	filesExpired := [testFileCount]string{}
	filesNonExpired := [testFileCount]string{}

	// 创建过期和未过期的文件
	for i := 0; i < testFileCount; i++ {
		num := ""
		if i > 0 {
			num = fmt.Sprintf("(%d)", i)
		}
		fileExpired := testName + "_" + expired.Format("060102150405") + "_" + strconv.Itoa(os.Getpid()) + num + ".txt"
		fileExpired = path.Join(testDir, fileExpired)
		os.WriteFile(fileExpired, nil, 0644)
		filesExpired[i] = fileExpired
		t.Logf("Created expired file: %s", filepath.Base(fileExpired))
	}
	for i := 0; i < testFileCount; i++ {
		num := ""
		if i > 0 {
			num = fmt.Sprintf("(%d)", i)
		}
		fileNonExpired := testName + "_" + now.Format("060102150405") + "_" + strconv.Itoa(os.Getpid()) + num + ".txt"
		fileNonExpired = path.Join(testDir, fileNonExpired)
		os.WriteFile(fileNonExpired, nil, 0644)
		filesNonExpired[i] = fileNonExpired
		t.Logf("Created non-expired file: %s", filepath.Base(fileNonExpired))
	}

	// 清理过期文件
	cleanExpireFiles(testDir, testName, expireTime)

	// 检查过期的文件是否已经删除，未过期的文件是否仍然存在
	for _, file := range filesExpired {
		_, err := os.Stat(file)
		if os.IsExist(err) {
			t.Errorf("Expired file not cleaned up: %s", filepath.Base(file))
		} else {
			t.Logf("Cleaned up expired file: %s", filepath.Base(file))
		}
	}
	for _, file := range filesNonExpired {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			t.Errorf("Non-expired file not kept: %s", filepath.Base(file))
		} else {
			t.Logf("Kept non-expired file: %s", filepath.Base(file))
		}
	}
}
