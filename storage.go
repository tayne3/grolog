// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package grolog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// 日志存储器
type groStorage struct {
	config       *Config
	file         *os.File
	out          *bufio.Writer
	lock         sync.Mutex
	err          error
	currFileNum  int
	currFileSize int64
}

// 创建新的存储器
func newStorage(config *Config) *groStorage {
	s := &groStorage{
		config: config,
	}

	if duration, err := time.ParseDuration(s.config.ExpireTime); err == nil && duration > 0 {
		cleanExpireFiles(s.config.FileDir, s.config.FileName, duration)
	}

	err := s.open(false)
	if err != nil {
		s.err = err
	}
	return s
}

// 停止存储器
func (s *groStorage) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.out == nil {
		return
	}

	s.closeFile()
}

// 刷新缓冲区
func (s *groStorage) Flush() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.out == nil {
		return
	}

	s.out.Flush()
	s.file.Sync()
}

// 错误
func (s *groStorage) Error() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.err
}

// 是否有效
func (s *groStorage) IsValid() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.out != nil
}

// 写入日志消息
func (s *groStorage) Write(b []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.out == nil {
		return
	}

	{
		size := int64(len(b))
		for size > 0 {
			available := s.config.MaxFileSize - s.currFileSize
			if size < available {
				s.out.Write(b)
				s.currFileSize += size
				break
			}
			s.out.Write(b[:available])
			s.currFileSize += size
			size -= available
			b = b[available:]
			if err := s.nextFile(); err != nil {
				return
			}
		}
	}

	if s.config.MaxWriteBuffer == 0 {
		s.out.Flush()
	}
}

// 写入日志消息
func (s *groStorage) WriteString(text string) {
	b := unsafe.Slice(unsafe.StringData(text), len(text))
	s.Write(b)
}

// 递归创建目录
func createNestedDirs(path string) error {
	dirs := strings.Split(path, string(filepath.Separator))
	for i := range dirs {
		subDir := strings.Join(dirs[:i+1], string(filepath.Separator))
		if _, err := os.Stat(subDir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err = os.MkdirAll(subDir, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

// 清理过期文件
func cleanExpireFiles(dir string, name string, expireTime time.Duration) {
	if expireTime <= 0 {
		return
	}
	expired := time.Now().Add(-expireTime)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (r error) {
		if err != nil {
			return
		}
		if info.IsDir() {
			return
		}

		// 使用正则表达式验证
		re := regexp.MustCompile(fmt.Sprintf(`^%s.+\.log$`, name))
		matches := re.FindStringSubmatch(info.Name())

		// 检查是否满足预期的日志文件名
		if matches == nil {
			return
		}

		if expired.Before(info.ModTime()) {
			os.Remove(path)
		}
		return
	})
}

// 获取文件大小
func (s *groStorage) getFileSize(file *os.File) (int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

// 打开日志文件
func (s *groStorage) open(clear bool) (err error) {
	var flag int
	if clear {
		flag = os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_APPEND
	}

	// 生成日志文件名 (Prefix_StartTime_PID(CurrFileNum).log)
	name := ""
	if s.config.EnableFileTime {
		num := ""
		if s.currFileNum > 0 {
			num = fmt.Sprintf("(%d)", s.currFileNum)
		}
		name = s.config.FileName + "_" + s.config.startTime.Format("060102150405") + "_" + strconv.Itoa(os.Getpid()) + num + ".log"
		name = filepath.Join(s.config.FileDir, name)
	} else {
		num := ""
		if s.currFileNum > 0 {
			num = fmt.Sprintf("(%d)", s.currFileNum)
		}
		name = filepath.Join(s.config.FileDir, s.config.FileName+num+".log")
	}

	// 目录不存在则创建
	if _, err := os.Stat(s.config.FileDir); os.IsNotExist(err) {
		createNestedDirs(s.config.FileDir)
	}
	// 文件不存在则创建
	if _, err := os.Stat(name); os.IsNotExist(err) {
		f, err := os.Create(name)
		if err != nil {
			fmt.Printf("create log file error, %v \n", err)
			return err
		}
		f.Close()
	} else {
		os.Chmod(name, 0644)
	}

	file, err := os.OpenFile(name, flag, 0644)
	if file == nil {
		fmt.Printf("open log file error, %v \n", err)
		return err
	}
	size, err := s.getFileSize(file)
	if err != nil {
		fmt.Printf("get log file size error, %v \n", err)
		return err
	}

	s.file = file
	s.out = bufio.NewWriterSize(s.file, s.config.MaxWriteBuffer)
	s.currFileSize = size
	return nil
}

// 关闭日志文件
func (s *groStorage) closeFile() {
	s.out.Flush()
	s.file.Close()
	s.out = nil
	s.file = nil
}

// 切换下一个文件序号
func (s *groStorage) nextFileNum() {
	s.currFileNum++
	if s.currFileNum >= s.config.MaxFileCount {
		s.currFileNum = 0
	}
}

// 切换下一个日志文件
func (s *groStorage) nextFile() error {
	s.closeFile()
	s.nextFileNum()
	if duration, err := time.ParseDuration(s.config.ExpireTime); err == nil && duration > 0 {
		cleanExpireFiles(s.config.FileDir, s.config.FileName, duration)
	}
	err := s.open(true)
	if err != nil {
		s.err = err
		return err
	}
	return nil
}
