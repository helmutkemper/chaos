package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"time"
)

// ImageInspect
//
// English:
//
//	Inspects the image and returns information type, ID, name, size, creation date, update date,
//	author, tags, etc
//
//	 Output:
//	   inspect: Image information
//	   err: Default error object
//
// Português:
//
//	Inspeciona a imagem e retorna informações tipo, ID, nome, tamanho, data de criação, data de
//	atualização, autor, tags, etc
//
//	 Saída:
//	  inspect: Informações da imagem
//	  err: Objeto de erro padrão
func (e *ContainerBuilder) ImageInspect() (inspect types.ImageInspect, err error) {
	if e.imageID == "" {
		e.imageID, err = e.ImageFindIdByName(e.GetImageName())
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	inspect, err = e.dockerSys.ImageInspect(e.imageID)
	if err != nil {
		util.TraceToLog()
		return
	}

	e.imageCreated, err = time.Parse(time.RFC3339Nano, inspect.Created)
	if err != nil {
		log.Printf("error: %v", err.Error())
		util.TraceToLog()
		return
	}

	e.imageInspected = true

	e.imageRepoTags = inspect.RepoTags
	e.imageRepoDigests = inspect.RepoDigests
	e.imageParent = inspect.Parent
	e.imageComment = inspect.Comment
	e.imageContainer = inspect.Container
	e.imageAuthor = inspect.Author
	e.imageArchitecture = inspect.Architecture
	e.imageVariant = inspect.Variant
	e.imageOs = inspect.Os
	e.imageOsVersion = inspect.OsVersion
	e.imageSize = inspect.Size
	e.imageVirtualSize = inspect.VirtualSize

	return
}
