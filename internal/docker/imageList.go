package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

// ImageList (English): List all images inside host server
//
// ImageList (Português): Lista todas as imagens do servidor hospedeiro
func (el *DockerSystem) ImageList() (
	list []types.ImageSummary,
	err error,
) {

	list, err = el.cli.ImageList(el.ctx, types.ImageListOptions{})
	return
}
