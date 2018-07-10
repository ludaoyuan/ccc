package common

import (
	"os"
	"path/filepath"
)

func CheckPathCreate(path string) error {
	dir := filepath.Dir(path)
	if CheckPath(dir) == false {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func CreatePath(path string) error {
	dir := filepath.Dir(path)

	if CheckPath(dir) == false {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除过期数据
func RemoveRegressionPath(path string) error {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
		} else {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
