package utils

import (
	"os"
	"path/filepath"
)

func DirExist(dir string) bool {
	s, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0755)
	return err == nil
}

// 删除文件（支持带通配符）
func DeletFiles(regString string) error {
	files, err := filepath.Glob(regString)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
