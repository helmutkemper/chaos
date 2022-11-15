package docker

import (
	"archive/tar"
	"bytes"
	"strings"
)

// ImageBuildPrepareFolderContext (English): Read the directory and prepare a .tar file
// header, used on backup tapes.
// Note: This function was made public to allow purposeful changes to the generated file,
// such as making it possible to add a Dockerfile file, for example.
//
// ImageBuildPrepareFolderContext (Português): Lê o diretório e prepara um header de
// arquivo .tar, usado em fitas de backup.
// Nota: Essa função foi deixada pública para permitir alterações propositais no arquivo
// gerado, como possibilitar adicionar um arquivo Dockerfile, por exemplo.
func (el DockerSystem) ImageBuildPrepareFolderContext(
	dirPath string,
) (
	file *bytes.Reader,
	err error,
) {

	var buf bytes.Buffer
	var tarWriter *tar.Writer
	tarWriter = tar.NewWriter(&buf)

	if strings.HasSuffix(dirPath, "/") == false {
		dirPath += "/"
	}

	err = el.ImageBuildPrepareFolderContextSupport(dirPath, dirPath, &buf, tarWriter)

	err = tarWriter.Close()
	if err != nil {
		return
	}

	file = bytes.NewReader(buf.Bytes())

	return
}
