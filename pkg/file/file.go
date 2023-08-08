package file

import (
	"fmt"
	"os"
)

// 写入文件，path:~/.config/poly/focus.yml
func CreateFileWithDir(path string, data []byte) error {
	fmt.Println(path)
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
