package iotmakerdocker

type ContainerNetworkDataList map[string]ContainerNetworkData

func (el ContainerNetworkDataList) GetGatewayByNetworkName(networkName string) (gateway string) {
	gateway = el[networkName].Gateway
	return
}

func (el ContainerNetworkDataList) GetIpAddressByNetworkName(networkName string) (ipAddress string) {
	ipAddress = el[networkName].IPAddress
	return
}

func (el ContainerNetworkDataList) GetEndpointIdByNetworkName(networkName string) (endpointID string) {
	endpointID = el[networkName].EndpointID
	return
}

func (el ContainerNetworkDataList) GetMacAddressByNetworkName(networkName string) (macAddress string) {
	macAddress = el[networkName].MacAddress
	return
}

func (el ContainerNetworkDataList) GetNetworkIdByNetworkName(networkName string) (networkID string) {
	networkID = el[networkName].NetworkID
	return
}
