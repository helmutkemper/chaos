package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func FileFindRecursively(fileName string) (filePath string, err error) {
	if _, err = os.Stat(fileName); os.IsNotExist(err) == false {
		filePath = fileName
		return
	}

	fileName = filepath.Base(fileName)
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == fileName {
				filePath = path
				return nil
			}

			return nil
		},
	)

	if filePath == "" {
		err = errors.New(fileName + ": file not found")
	}

	return
}

func FileFindContainsRecursively(fileName string) (filePath string, err error) {
	if _, err = os.Stat(fileName); os.IsNotExist(err) == false {
		filePath = fileName
		return
	}

	fileName = filepath.Base(fileName)
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(info.Name(), fileName) == true {
				filePath = path
				return nil
			}

			return nil
		},
	)

	if filePath == "" {
		err = errors.New(fileName + ": file not found")
	}

	return
}

func FileFindHasPrefixRecursively(fileName string) (filePath string, err error) {
	if _, err = os.Stat(fileName); os.IsNotExist(err) == false {
		filePath = fileName
		return
	}

	fileName = filepath.Base(fileName)
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasPrefix(info.Name(), fileName) == true {
				filePath = path
				return nil
			}

			return nil
		},
	)

	if filePath == "" {
		err = errors.New(fileName + ": file not found")
	}

	return
}

func FileFindContainsRecursivelyFullPath(fileName string) (filePath string, err error) {
	filePath, err = FileFindContainsRecursively(fileName)
	if err != nil {
		return
	}

	filePath, err = filepath.Abs(filePath)
	return
}

func FileFindHasPrefixRecursivelyFullPath(fileName string) (filePath string, err error) {
	filePath, err = FileFindHasPrefixRecursively(fileName)
	if err != nil {
		return
	}

	filePath, err = filepath.Abs(filePath)
	return
}

func FileFindRecursivelyFullPath(fileName string) (filePath string, err error) {
	filePath, err = FileFindRecursively(fileName)
	if err != nil {
		return
	}

	filePath, err = filepath.Abs(filePath)
	return
}
