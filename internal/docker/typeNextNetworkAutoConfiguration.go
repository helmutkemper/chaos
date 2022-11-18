package iotmakerdocker

import (
	"github.com/docker/docker/api/types/network"
)

type NextNetworkAutoConfiguration struct {
	ip      IPv4Generator
	id      string
	name    string
	gateway string
	subnet  string
	err     error
}

// init a network for new container
// nInit("test", 10, 0, 0, 1)
// before use this function, use whaleAquarium.Docker.NetworkCreate("test")
func (el *NextNetworkAutoConfiguration) Init(id, name, gateway, subnet string) {
	el.id = id
	el.name = name
	el.gateway = gateway
	el.subnet = subnet
	el.err = el.ip.InitWithString(gateway, subnet)
}

func (el *NextNetworkAutoConfiguration) GetNext() (IP string, networkConfig *network.NetworkingConfig, err error) {
	IP = el.ip.String()
	el.err = el.ip.IncCurrentIP()
	return IP,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				el.name: {
					NetworkID: el.id,
					Gateway:   el.gateway,
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: IP,
					},
					IPAddress: IP,
				},
			},
		},
		el.err
}

func (el *NextNetworkAutoConfiguration) GetCurrentIpAddress() (IP string, err error) {
	return el.ip.String(), el.err
}
