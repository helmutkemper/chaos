package docker

import (
	"github.com/helmutkemper/util"
)

// GetNetworkIPV4
//
// English:
//
//	Return the IPV4 from the docker network
//
//	 Output:
//	   IPV4: network address IPV4
//
// Português:
//
//	Retorno o IPV4 da rede do docker
//
//	 Saída:
//	   IPV4: endereço IPV4 da rede
func (e *ContainerBuilder) GetNetworkIPV4() (IPV4 string) {
	var err error
	var inspect ContainerInspect

	inspect, err = e.ContainerInspect()
	if err != nil {
		util.TraceToLog()
		return
	}

	IPV4 = inspect.Network.IPAddress
	return
}
