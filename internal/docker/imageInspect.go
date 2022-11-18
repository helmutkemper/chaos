package iotmakerdocker

import "github.com/docker/docker/api/types"

// ImageInspect (English): Inspect image by ID
//
// ImageInspect (Português): Inspeciona a imagem por ID
func (el *DockerSystem) ImageInspect(id string) (inspect types.ImageInspect, err error) {

	inspect, _, err = el.cli.ImageInspectWithRaw(el.ctx, id)
	if err != nil {
		return
	}

	return
}
