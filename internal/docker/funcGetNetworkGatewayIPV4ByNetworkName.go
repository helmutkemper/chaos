package docker

import (
	"errors"
	"github.com/helmutkemper/util"
)

// GetNetworkGatewayIPV4ByNetworkName
//
// English:
//
//	If the container is connected to more than one network, this function returns the gateway of the
//	chosen network.
//
//	 Input:
//	   networkName: name of the network
//
//	 Output:
//	   IPV4: address of the gateway of the network
//	   err: standard object error
//
// Note:
//
//   - The default docker network is named "bridge"
//
// Português:
//
//	Caso o container esteja ligado em mais de uma rede, esta função devolve o gateway da rede
//	escolhida.
//
//	 Entrada:
//	   networkName: nome da rede
//
//	 Saída:
//	   IPV4: endereço do gateway da rede
//	   err: objeto de erro padrão
//
// Nota:
//
//   - A rede padrão do docker tem o nome "bridge"
func (e *ContainerBuilder) GetNetworkGatewayIPV4ByNetworkName(networkName string) (IPV4 string, err error) {
	var found bool
	var inspect ContainerInspect

	inspect, err = e.ContainerInspect()
	if err != nil {
		util.TraceToLog()
		return
	}

	_, found = inspect.Network.Networks[networkName]
	if found == false {
		util.TraceToLog()
		err = errors.New("network name not found")
		return
	}

	IPV4 = inspect.Network.Networks[networkName].Gateway
	return
}
