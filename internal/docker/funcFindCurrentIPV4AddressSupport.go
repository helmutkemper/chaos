package docker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"strings"
)

// findCurrentIPV4AddressSupport
//
// English:
//
//	Support function for FindCurrentIpAddress()
//
//	 Input:
//	   networkID: Docker's network ID
//
//	 Output:
//	   IP: network IPV4 address
//	   err: standard error object
//
// Português: função de apoio a FindCurrentIpAddress()
//
//	Entrada:
//	  networkID: ID da rede docker
//
//	Saída:
//	  IP: endereço IPV4 da rede
//	  err: objeto de erro padrão
func (e *ContainerBuilder) findCurrentIPV4AddressSupport(networkID string) (IP string, err error) {
	var res types.NetworkResource
	res, err = e.dockerSys.NetworkInspect(networkID)
	if err != nil {
		util.TraceToLog()
		return
	}

	var pass = false
	for containerID, networkData := range res.Containers {
		if containerID == e.containerID && networkData.Name == e.containerName {
			pass = true
			var parts = strings.Split(networkData.IPv4Address, "/")
			IP = parts[0]
			return
		}
	}

	if pass == false {
		util.TraceToLog()
		err = errors.New("container not found on bridge network")
		return
	}

	return
}
