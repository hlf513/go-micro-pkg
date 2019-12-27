package utils

import (
	"os"
	"path/filepath"
)

// PathExist 判断 path 是否存在
func PathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// CurrentDir 获取当前目录
func CurrentDir() string {
	var sp = string(filepath.Separator)
	dirPath, _ := filepath.Abs(filepath.Dir(filepath.Join("."+sp, sp)))
	return dirPath
}
