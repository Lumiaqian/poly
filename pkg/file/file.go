package file

import (
	"os"
)

// 写入文件，path:./ fileName:focus.yml
func CreateFileWithDir(path, fileName string, data []byte) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(path + fileName)
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
