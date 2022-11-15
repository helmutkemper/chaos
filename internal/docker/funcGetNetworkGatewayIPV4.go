package docker

import (
	"github.com/helmutkemper/util"
)

// GetNetworkGatewayIPV4
//
// English:
//
//	Returns the gateway from the network to the IPV4 network
//
//	 Output:
//	   IPV4: IPV4 address of the gateway
//
// Português:
//
//	Retorna o gateway da rede para rede IPV4
//
//	 Saída:
//	   IPV4: endereço IPV4 do gateway
func (e *ContainerBuilder) GetNetworkGatewayIPV4() (IPV4 string) {
	var err error
	var inspect ContainerInspect

	inspect, err = e.ContainerInspect()
	if err != nil {
		util.TraceToLog()
		return
	}

	IPV4 = inspect.Network.Gateway
	return
}
