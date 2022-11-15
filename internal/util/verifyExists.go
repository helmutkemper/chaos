package util

import "os"

func VerifyExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() == true {
		return true
	}

	return true //!info.IsDir()
	//old
	//_, err := os.Stat(path)
	//if os.IsNotExist(err) {
	//	return false
	//}
	//return true
}
