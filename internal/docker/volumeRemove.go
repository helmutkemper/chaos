package docker

func (el *DockerSystem) VolumeRemove(
	name string,
) (
	err error,
) {

	err = el.cli.VolumeRemove(el.ctx, name, false)
	return
}
