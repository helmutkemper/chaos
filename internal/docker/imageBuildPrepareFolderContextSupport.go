package docker

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

func (el DockerSystem) ImageBuildPrepareFolderContextSupport(
	dirPath,
	toRemoveInsideTarFilePathList string,
	buf *bytes.Buffer,
	tarWriter *tar.Writer,
) (
	err error,
) {

	var dirContent []os.FileInfo
	var tarHeader *tar.Header
	var fileData []byte
	var filePath string

	if strings.HasSuffix(dirPath, "/") == false {
		dirPath += "/"
	}

	dirContent, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, folderItem := range dirContent {
		filePath = dirPath + folderItem.Name()

		if folderItem.IsDir() == true {
			err = el.ImageBuildPrepareFolderContextSupport(filePath, toRemoveInsideTarFilePathList, buf, tarWriter)
			if err != nil {
				return
			}
		} else {
			fileData, err = ioutil.ReadFile(filePath)
			if err != nil {
				return
			}

			xxx := strings.Replace(filePath, toRemoveInsideTarFilePathList, "", 1)
			_ = xxx

			tarHeader = &tar.Header{
				Name: strings.Replace(filePath, toRemoveInsideTarFilePathList, "", 1),
				Mode: 0600,
				Size: folderItem.Size(),
			}

			if err = tarWriter.WriteHeader(tarHeader); err != nil {
				return
			}
			if _, err = tarWriter.Write(fileData); err != nil {
				return
			}
		}
	}

	return
}
