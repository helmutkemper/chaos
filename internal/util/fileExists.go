package util

import "os"

func VerifyFileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() == true {
		return true
	}

	return true //!info.IsDir()
}
