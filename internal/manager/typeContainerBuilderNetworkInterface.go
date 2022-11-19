package manager

import (
	networkTypes "github.com/docker/docker/api/types/network"
)

type NetworkInterface interface {
	Init() (err error)
	GetConfiguration() (IP string, networkConfiguration *networkTypes.NetworkingConfig, err error)
	NetworkCreate(name, subnet, gateway string) (err error)
	GetNetworkName() (name string)
	GetNetworkID() (ID string)
	Remove() (err error)
}
