package file

import (
	"os"
	"path/filepath"
)

// 写入文件，path:~/.config/poly/focus.yml
func CreateFileWithDir(path string, data []byte) error {
	// 确保文件所在的目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
