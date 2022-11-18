package builder

import "github.com/docker/docker/api/types"

// ContainerNetworkInspect (English): Returns athe network data from a container
//
//	id: string container id
//
// ContainerNetworkInspect (Português): Retorna os dados da rede do container
//
//	id: string container id
func (el *DockerSystem) ContainerNetworkInspect(
	id string,
) (
	netDataList ContainerNetworkDataList,
	err error,
) {

	var insp types.ContainerJSON

	netDataList = make(map[string]ContainerNetworkData)

	insp, err = el.ContainerInspect(id)
	if err != nil {
		return
	}

	containerNetworks := (*insp.NetworkSettings).Networks
	for k, v := range containerNetworks {
		netDataList[k] = ContainerNetworkData{
			Gateway:    v.Gateway,
			IPAddress:  v.IPAddress,
			EndpointID: v.EndpointID,
			NetworkID:  v.NetworkID,
			MacAddress: v.MacAddress,
		}
	}

	return
}
