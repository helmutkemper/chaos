package docker

import (
	"github.com/helmutkemper/util"
)

// ImageListExposedVolumes
//
// English:
//
//	Lists all volumes defined in the image.
//
//	 Output:
//	   list: List of exposed volumes
//	   err: Standard error object
//
// Note:
//
//   - Use the AddFileOrFolderToLinkBetweenComputerHostAndContainer() function to link folders and
//     files between the host computer and the container
//
// Português:
//
//	Lista todos os volumes definidos na imagem.
//
//	 Saída:
//		 list: Lista de volumes expostos
//	   err: Objeto de erro padrão
//
// Nota:
//
//   - Use a função AddFileOrFolderToLinkBetweenComputerHostAndContainer() para vincular pastas e
//     arquivos entre o computador hospedeiro e o container
func (e *ContainerBuilder) ImageListExposedVolumes() (list []string, err error) {

	list, err = e.dockerSys.ImageListExposedVolumes(e.imageID)
	if err != nil {
		util.TraceToLog()
	}
	return
}
