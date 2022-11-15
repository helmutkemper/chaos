package util

import "os"

func VerifyDirExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
