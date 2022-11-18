package builder

import (
	"fmt"
)

// NewNetworkWithHighAddress (English): Create a network with gateway address.0.0.1
// and subnet address.0.0.0/subnetMask
//
//	networkName: network name
//	address: upper part of the address between 1 and 255
//	subnetMask: subnet mask. Examples: 2, 4, 8, 16, 32
//
// NewNetworkWithHighAddress (Português): Cria uma rede com gateway address.0.0.1
// e subnet address.0.0.0/subnetMask
//
//	networkName: nome da rede
//	address: parte alta do endereço entre 1 e 255
//	subnetMask: máscara de rede. Exemplos: 2, 4, 8, 16, 32
func NewNetworkWithHighAddress(
	networkName string,
	address byte,
	subnetMask byte,
) (
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
		fmt.Sprintf("%d.0.0.0/%d", address, subnetMask),
		fmt.Sprintf("%d.0.0.1", address),
	)

	return
}
