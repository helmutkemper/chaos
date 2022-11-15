package docker

func (el *DockerSystem) NetworkListLocalCreated() (
	list map[string]string,
) {

	return el.networkId
}
