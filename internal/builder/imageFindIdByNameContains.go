package builder

import (
	"errors"
	"github.com/docker/docker/api/types"
	"strings"
)

// ImageFindIdByName (English): Find image whose part of the name contains the search
// term
//
//	name: search term
//
// ImageFindIdByName (Português): Procura uma imagem cujo o nome contém o termo procurado
//
//	name: termo procurado
func (el *DockerSystem) ImageFindIdByNameContains(
	name string,
) (
	list []NameAndId,
	err error,
) {

	var listTmp []types.ImageSummary
	list = make([]NameAndId, 0)

	listTmp, err = el.ImageList()
	if err != nil {
		return
	}

	if len(el.imageId) == 0 {
		el.imageId = make(map[string]string)
	}

	for _, data := range listTmp {
		if len(data.RepoTags) == 0 {
			continue
		}

		var tag = data.RepoTags[0]
		el.imageId[tag] = data.ID
		if strings.Contains(tag, name) == true {
			list = append(list, NameAndId{
				ID:   data.ID,
				Name: tag,
			})
		}
	}

	if len(list) == 0 {
		err = errors.New("image name not found")
	}

	return
}
