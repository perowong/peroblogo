package utils

import (
	"fmt"
	"os"
)

func ReadFile(filename string) (data []byte, err error) {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("file %s is not exist", filename)
		}
		return data, err
	}

	if fileinfo.IsDir() {
		err = fmt.Errorf("file %s is dir", filename)
		return data, err
	}

	return os.ReadFile(filename)
}
