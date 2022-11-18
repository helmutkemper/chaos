package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"io"
)

func (el *DockerSystem) ImagePull(
	name string,
	channel *chan ContainerPullStatusSendToChannel,
) (
	imageId string,
	imageName string,
	err error,
) {

	var reader io.Reader

	//esse valor é trocado no final do download
	imageName = name

	reader, err = el.cli.ImagePull(el.ctx, name, types.ImagePullOptions{})
	if err != nil {
		return
	}

	if len(el.imageId) == 0 {
		el.imageId = make(map[string]string)
	}

	el.imageId[name] = ""
	var successfully bool
	successfully, err = el.processBuildAndPullReaders(&reader, channel)
	if successfully == false || err != nil {
		if err != nil {
			return
		}

		err = errors.New("image pull error")
	}

	imageId, err = el.ImageFindIdByName(name)
	if err != nil {
		return
	}

	return
}
