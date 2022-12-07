package manager

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/builder"
)

type Network struct {
	manager *Manager
}

func (el *Network) New(manager *Manager) {
	el.manager = manager
	el.manager.network = new(dockerNetwork)
}

// NetworkCreate
//
// Create a docker network to be used in the chaos test
//
//	Input:
//	  name: network name
//	  subnet: subnet value. eg. 10.0.0.0/16
//	  gateway: gateway value. eg. "10.0.0.1
//
//	Notes:
//	  * If there is already a network with the same name and the same configuration, nothing will be done;
//	  * If a network with the same name and different configuration already exists, the network will be deleted and a new network created.
func (el *Network) NetworkCreate(name, subnet, gateway string) (err error) {
	el.manager.network.networkName = name

	var networkList []types.NetworkResource
	networkList, err = el.manager.DockerSys[0].NetworkList()
	if err != nil {
		return
	}

	for _, networkData := range networkList {
		if networkData.Name == name {
			el.manager.network.networkID = networkData.ID

			var data types.NetworkResource
			if data, err = el.manager.DockerSys[0].NetworkInspect(networkData.ID); err != nil {
				err = fmt.Errorf("network.NetworkCreate().NetworkInspect().error: %v", err)
				return
			}
			if data.IPAM.Config[0].Subnet != subnet || data.IPAM.Config[0].Gateway != gateway {
				if err = el.manager.DockerSys[0].NetworkRemove(networkData.ID); err != nil {
					err = fmt.Errorf("network.NetworkCreate().NetworkRemove().error: %v", err)
					return
				}

				break
			}

			el.manager.network.generator = el.manager.DockerSys[0].NetworkGetGenerator(name)
			return
		}
	}

	el.manager.network.networkID, el.manager.network.generator, err = el.manager.DockerSys[0].NetworkCreate(name, builder.KNetworkDriveBridge, "local", subnet, gateway)
	return
}
