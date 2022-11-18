package builder

func (el *DockerSystem) NetworkDisconnect(
	networkID,
	containerID string,
	force bool,
) (
	err error,
) {

	return el.cli.NetworkDisconnect(el.ctx, networkID, containerID, force)
}
