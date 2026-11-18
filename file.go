package gokit

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

/**** 文件操作 ****/

// 文件是否存在（不含权限）
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// 按行读取文件。 单行最大长度为 16MB，超过则返回 bufio.ErrTooLong
func ReadFileByLine(path string, fn func(string) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// 设置初始缓冲区为 1MB，最大限制为 16MB（可按需调整）
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 16*1024*1024)

	for scanner.Scan() {
		if err = fn(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// 读取CSV （fn内的[]string为复用结构，禁止直接保存或传递）
func ReadCSV(path string, fn func([]string) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.ReuseRecord = true // 复用同一底层数组，减少分配

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // 读取完毕
		}
		if err != nil {
			return err
		}
		if err = fn(record); err != nil {
			return err
		}
	}
	return nil
}
