package docker

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/helmutkemper/util"
	"path/filepath"
)

// AddFileOrFolderToLinkBetweenComputerHostAndContainer
//
// English:
//
//	Links a file or folder between the computer host and the container.
//
//	 Input:
//	   computerHostPath:    Path of the file or folder inside the host computer.
//	   insideContainerPath: Path inside the container.
//
//	 Output:
//	   err: Default error object.
//
// Português:
//
//	Vincula um arquivo ou pasta entre o computador e o container.
//
//	 Entrada:
//	   computerHostPath:    Caminho do arquivo ou pasta no computador hospedeiro.
//	   insideContainerPath: Caminho dentro do container.
//
//	 Output:
//	   err: Objeto de erro padrão.
func (e *ContainerBuilder) AddFileOrFolderToLinkBetweenComputerHostAndContainer(computerHostPath, insideContainerPath string) (err error) {

	if e.volumes == nil {
		e.volumes = make([]mount.Mount, 0)
	}

	computerHostPath, err = filepath.Abs(computerHostPath)
	if err != nil {
		util.TraceToLog()
		return
	}

	e.volumes = append(
		e.volumes,
		mount.Mount{
			// bind - is the type for mounting host dir (real folder inside computer where this code work)
			Type: KVolumeMountTypeBindString,
			// path inside host machine
			Source: computerHostPath,
			// path inside image
			Target: insideContainerPath,
		},
	)

	return
}
