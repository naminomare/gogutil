package fileio

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

var (
	// ErrOverMaxLoop maxLoopを超えてもファイルが存在し続けた場合のエラー
	ErrOverMaxLoop = errors.New("maxLoopを超えてもファイルが存在し続けた場合のエラー")
)

// IsExist ファイルが存在するかチェック
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetNonExistFileName pathが存在した場合に、path0, path1のようなファイル名を返す
func GetNonExistFileName(path string, maxLoop int) (string, error) {
	if !IsExist(path) {
		return path, nil
	}

	// 存在した場合
	ext := filepath.Ext(path)
	filenamebase := path[0 : len(path)-len(ext)]

	for i := 0; i < maxLoop; i++ {
		filename := filenamebase + strconv.Itoa(i) + ext
		if IsExist(filename) {
			continue
		}
		return filename, nil
	}
	return "", ErrOverMaxLoop
}
