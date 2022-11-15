package docker

import (
	"github.com/docker/docker/api/types/network"
	"github.com/helmutkemper/util"
)

// NetworkChangeIp
//
// English:
//
//	Change the IP address of the container, to the next IP in the docker network manager list
//
//	 Output:
//	   err: Default object error from golang
//
// Português:
//
//	Troca o endereço IP do container, para o próximo IP da lista do gerenciador de rede docker
//
//	 Saída:
//	   err: Objeto padrão de erro golang
func (e *ContainerBuilder) NetworkChangeIp() (err error) {
	var networkID = e.network.GetNetworkID()
	err = e.dockerSys.NetworkDisconnect(networkID, e.containerID, false)
	if err != nil {
		util.TraceToLog()
		return
	}

	var netConfig *network.NetworkingConfig
	e.IPV4Address, netConfig, err = e.network.GetConfiguration()
	if err != nil {
		util.TraceToLog()
		return
	}

	err = e.dockerSys.NetworkConnect(networkID, e.containerID, netConfig.EndpointsConfig[e.network.GetNetworkName()])
	if err != nil {
		util.TraceToLog()
		return
	}

	return
}
