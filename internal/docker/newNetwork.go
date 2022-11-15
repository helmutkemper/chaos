package docker

// NewNetwork (English): Create a network with gateway 10.0.0.1 and subnet 10.0.0.0/16
//
//	networkName: network name
//
// NewNetwork (Português): Cria uma rede com gateway 10.0.0.1 e subnet 10.0.0.0/16
func NewNetwork(networkName string) (
	networkId string,
	networkAutoConfiguration *NextNetworkAutoConfiguration,
	err error,
) {

	var dockerSys = DockerSystem{}
	err = dockerSys.Init()
	if err != nil {
		return
	}

	networkId, networkAutoConfiguration, err = dockerSys.NetworkCreate(
		networkName,
		KNetworkDriveBridge,
		"local",
		"10.0.0.0/16",
		"10.0.0.1",
	)

	return
}
