package util

import "path/filepath"

func FileGetAbsolutePath(pathFile string) (error, string) {
	var err error
	var absolutePath string
	absolutePath, err = filepath.Abs(pathFile)

	return err, absolutePath
}
