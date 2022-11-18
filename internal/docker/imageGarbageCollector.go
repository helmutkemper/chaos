package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

// ImageGarbageCollector (English): Remove temporary images, where first tag is
// "<none>:<none>"
//
// ImageGarbageCollector (Português): Remove imagens temporárias, onde a primeira tag é
// "<none>:<none>"
func (el *DockerSystem) ImageGarbageCollector() (err error) {
	var list []types.ImageSummary
	list, err = el.ImageList()
	for _, image := range list {
		if len(image.RepoTags) > 0 {
			if image.RepoTags[0] == "<none>:<none>" {
				err = el.ImageRemove(image.ID, true, true)
				if err != nil {
					return
				}
			}
		}
	}

	return
}
