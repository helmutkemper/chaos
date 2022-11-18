package iotmakerdocker

func (el *DockerSystem) NetworkRemove(
	id string,
) (
	err error,
) {

	return el.cli.NetworkRemove(el.ctx, id)
}
