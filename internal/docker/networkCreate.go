package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// NetworkCreate create network
//
//	name:    string       Ex.: containerNetwork
//	drive:   NetworkDrive Ex.: KNetworkDriveBridge
//	scope:   string       Ex.: local
//	subnet:  string       Ex.: 10.0.0.0/16 (note: use base 10)
//	gateway: string       Ex.: 10.0.0.1    (note: use base 10)
func (el *DockerSystem) NetworkCreate(
	name string,
	drive NetworkDrive,
	scope,
	subnet,
	gateway string,
) (
	id string,
	networkGenerator *NextNetworkAutoConfiguration,
	err error,
) {

	//todo: se já tem uma rede, ajustar o ip automático para o próximo endereço
	var resp types.NetworkCreateResponse
	var insp types.NetworkResource

	networkGenerator = &NextNetworkAutoConfiguration{}

	if len(el.networkId) == 0 {
		el.networkId = make(map[string]string)
	}

	if len(el.networkGenerator) == 0 {
		el.networkGenerator = make(map[string]*NextNetworkAutoConfiguration)
	}

	id, _ = el.NetworkFindIdByName(name)
	if id != "" {

		insp, err = el.cli.NetworkInspect(
			el.ctx,
			id,
			types.NetworkInspectOptions{
				Scope:   scope,
				Verbose: false,
			},
		)
		if err != nil {
			return
		}
		pass := false
		for _, v := range insp.IPAM.Config {
			if v.Gateway == gateway && v.Subnet == subnet {
				pass = true
				break
			}
		}

		if pass == true {

			var res types.NetworkResource
			res, err = el.cli.NetworkInspect(el.ctx, name, types.NetworkInspectOptions{
				Scope:   scope,
				Verbose: false,
			})

			if err != nil {
				return
			}

			var biggestIP = "0.0.0.0"
			for _, containerNetwork := range res.Containers {
				biggestIP, err = el.networkGetTheBiggestAddress(biggestIP, containerNetwork.IPv4Address)
			}

			networkGenerator.Init(
				res.ID,
				name,
				gateway,
				subnet,
			)

			el.networkId[name] = res.ID
			id = res.ID

			el.networkGenerator[name] = networkGenerator

			return
		}

		err = errors.New("there is a network with this name")
		return
	}

	resp, err = el.cli.NetworkCreate(el.ctx, name, types.NetworkCreate{
		//CheckDuplicate: false,
		Driver: drive.String(),
		Scope:  scope,
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet:  subnet,
					Gateway: gateway,
				},
			},
		},
		Attachable: true,
		Labels: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return
	}

	networkGenerator.Init(
		resp.ID,
		name,
		gateway,
		subnet,
	)

	el.networkId[name] = resp.ID
	id = resp.ID

	el.networkGenerator[name] = networkGenerator

	return
}
