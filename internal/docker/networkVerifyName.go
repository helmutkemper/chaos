package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

// verify if network name exists
func (el *DockerSystem) NetworkVerifyName(
	name string,
) (
	exists bool,
	err error,
) {

	var resp []types.NetworkResource
	resp, err = el.cli.NetworkList(el.ctx, types.NetworkListOptions{})
	if err != nil {
		return false, err
	}

	for _, v := range resp {
		if v.Name == name {
			return true, nil
		}
	}

	return false, nil
}
