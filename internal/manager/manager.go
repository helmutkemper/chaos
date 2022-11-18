package manager

import "github.com/docker/docker/api/types"

type Manager struct {
}

type Network struct {
	manager *Manager
}

func (el *Network) New(manager *Manager) {
	el.manager = manager
}

func (el *Network) NetworkCreate(name, subnet, gateway string) (err error) {
	el.networkName = name

	var networkList []types.NetworkResource
	networkList, err = el.dockerSys.NetworkList()
	if err != nil {
		return
	}

	for _, networkData := range networkList {
		if networkData.Name == name {
			el.networkID = networkData.ID
			el.generator = el.dockerSys.NetworkGetGenerator(name)
			return
		}
	}

	el.networkID, el.generator, err = el.dockerSys.NetworkCreate(name, docker.KNetworkDriveBridge, "local", subnet, gateway)
	return
}
