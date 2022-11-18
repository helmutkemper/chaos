package iotmakerdocker

func (el *DockerSystem) NetworkRemoveByName(
	name string,
) (
	err error,
) {

	var id string
	id, err = el.NetworkFindIdByName(name)
	if err != nil {
		return err
	}

	return el.cli.NetworkRemove(el.ctx, id)
}
