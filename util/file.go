package util

import (
	"os"
	"path/filepath"
)

// 保存文件
func SaveFile(path, fileName string, data []byte) error {
	// 创建目录
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	// 创建文件
	filePath := filepath.Join(path, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 读取文件
func ReadFile(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return content
}

// 检查是否存在
func CheckExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err) || err == nil
}
