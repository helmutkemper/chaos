package docker

import (
	isolatedNetwork "github.com/helmutkemper/iotmaker.docker.builder.network.interface"
)

// GetNetworkInterface
//
// English:
//
//	Returns the object defined for the network control
//
//	 Output:
//	   network: Object pointer used to configure the network
//
// Português:
//
//	Retorna o objeto definido para o controle da rede
//
//	 Saída:
//	   network: Ponteiro do objeto usado para configurar a rede
func (e *ContainerBuilder) GetNetworkInterface() (network isolatedNetwork.ContainerBuilderNetworkInterface) {
	return e.network
}
