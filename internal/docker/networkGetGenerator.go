package docker

func (el *DockerSystem) NetworkGetGenerator(
	name string,
) (
	configuration *NextNetworkAutoConfiguration,
) {

	return el.networkGenerator[name]
}
