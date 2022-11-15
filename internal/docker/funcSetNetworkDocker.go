package docker

import (
	isolatedNetwork "github.com/helmutkemper/iotmaker.docker.builder.network.interface"
)

// SetNetworkDocker
//
// English:
//
//	Sets the docker network manager pointer
//
//	 Input:
//	   network: pointer to the network manager object.
//
// Note:
//
//   - Compatible with dockerBuilderNetwork.ContainerBuilderNetwork{} object
//
// Português:
//
//	Define o ponteiro do gerenciador de rede docker
//
//	 Entrada:
//	   network: ponteiro para o objeto gerenciador de rede.
//
// Nota:
//
//   - Compatível com o objeto dockerBuilderNetwork.ContainerBuilderNetwork{}
func (e *ContainerBuilder) SetNetworkDocker(network isolatedNetwork.ContainerBuilderNetworkInterface) {
	e.network = network
}
