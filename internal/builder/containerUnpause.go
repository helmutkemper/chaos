package builder

func (el *DockerSystem) ContainerUnpause(
	id string,
) (
	err error,
) {

	return el.cli.ContainerUnpause(el.ctx, id)
}
