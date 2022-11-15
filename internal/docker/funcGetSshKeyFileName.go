package docker

import (
	"io/fs"
	"io/ioutil"
	"math"
	"path"
	"strings"
)

// GetSshKeyFileName
//
// English:
//
//	Returns the name of the last generated ssh key.
//
// Português:
//
//	Retorna o nome da chave ssh gerada por último.
func (e *ContainerBuilder) GetSshKeyFileName(dir string) (fileName string, err error) {
	var folderPath = path.Join(dir, ".ssh")

	var minDate = int64(math.MaxInt64)

	var files []fs.FileInfo
	files, err = ioutil.ReadDir(folderPath)

	for _, file := range files {
		var name = file.Name()
		var date = file.ModTime().UnixNano()

		if file.IsDir() == true {
			continue
		}

		if strings.HasPrefix(name, "id_") == true && strings.HasSuffix(name, ".pub") == false {
			if minDate >= date {
				minDate = date
				fileName = name
			}
		}
	}

	return
}
