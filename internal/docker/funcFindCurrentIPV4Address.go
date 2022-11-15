package docker

import (
	"github.com/helmutkemper/util"
)

// FindCurrentIPV4Address
//
// English:
//
//	Inspects the docker's network and returns the current IP of the container
//
//	 Output:
//	   IP: container IP address IPV4
//	   err: standard error object
//
// Português:
//
//	Inspeciona a rede do docker e devolve o IP atual do container
//
//	 Saída:
//	   IP: endereço IP do container IPV4
//	   err: objeto de erro padrão
func (e *ContainerBuilder) FindCurrentIPV4Address() (IP string, err error) {
	var id string
	if e.network == nil {
		id, err = e.dockerSys.NetworkFindIdByName("bridge")
		if err != nil {
			util.TraceToLog()
			return
		}
		IP, err = e.findCurrentIPV4AddressSupport(id)
		if err != nil {
			util.TraceToLog()
		}
	} else {
		IP, err = e.findCurrentIPV4AddressSupport(e.network.GetNetworkID())
		if err != nil {
			util.TraceToLog()
		}
	}

	return
}
