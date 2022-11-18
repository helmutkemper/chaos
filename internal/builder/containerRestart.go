package builder

func (el *DockerSystem) ContainerRestart(
	id string,
) (
	err error,
) {

	return el.cli.ContainerRestart(el.ctx, id, nil)
}
