package manager

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/builder"
)

type Network struct {
	manager *Manager

	generator   *builder.NextNetworkAutoConfiguration
	networkID   string
	networkName string
}

func (el *Network) New(manager *Manager) {
	el.manager = manager
}

// NetworkCreate
//
// Create a docker network to be used in the chaos test
//
//	Input:
//	  name: network name
//	  subnet: subnet value. eg. 10.0.0.0/16
//	  gateway: gateway value. eg. "10.0.0.1
func (el *Network) NetworkCreate(name, subnet, gateway string) (err error) {
	el.networkName = name

	var networkList []types.NetworkResource
	networkList, err = el.manager.DockerSys.NetworkList()
	if err != nil {
		return
	}

	for _, networkData := range networkList {
		if networkData.Name == name {
			el.networkID = networkData.ID

			var data types.NetworkResource
			if data, err = el.manager.DockerSys.NetworkInspect(networkData.ID); err != nil {
				err = fmt.Errorf("network.NetworkCreate().NetworkInspect().error: %v", err)
				return
			}
			if data.IPAM.Config[0].Subnet != subnet || data.IPAM.Config[0].Gateway != gateway {
				if err = el.manager.DockerSys.NetworkRemove(networkData.ID); err != nil {
					err = fmt.Errorf("network.NetworkCreate().NetworkRemove().error: %v", err)
					return
				}

				break
			}

			el.generator = el.manager.DockerSys.NetworkGetGenerator(name)
			return
		}
	}

	el.networkID, el.generator, err = el.manager.DockerSys.NetworkCreate(name, builder.KNetworkDriveBridge, "local", subnet, gateway)
	return
}
