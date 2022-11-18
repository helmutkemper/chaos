package builder

func (el *DockerSystem) ContainerPause(
	id string,
) (
	err error,
) {

	return el.cli.ContainerPause(el.ctx, id)
}
