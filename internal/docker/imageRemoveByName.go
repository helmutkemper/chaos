package iotmakerdocker

func (el *DockerSystem) ImageRemoveByName(
	name string,
	force,
	pruneChildren bool,
) (
	err error,
) {

	var id string

	id, err = el.ImageFindIdByName(name)

	if err != nil {
		return err
	}

	err = el.ImageRemove(id, force, pruneChildren)

	return
}
