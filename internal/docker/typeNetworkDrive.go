package iotmakerdocker

type NetworkDrive int

func (el NetworkDrive) String() string {
	return networksDrivers[el]
}

const (
	//KNetworkDriveBridge (English): The default network driver.
	//
	//If you don’t specify a driver, this is the type of network you are creating. Bridge networks are usually used
	//when your applications run in standalone containers that need to communicate.
	//
	//See bridge networks https://docs.docker.com/network/bridge/.
	KNetworkDriveBridge NetworkDrive = iota

	//KNetworkDriveHost (English): For standalone containers, remove network isolation between the container and the
	//Docker host, and use the host’s networking directly.
	//
	//Host is only available for swarm services on Docker 17.06 and higher.
	//
	//See use the host network. https://docs.docker.com/network/host/
	KNetworkDriveHost

	//KNetworkDriveOverlay (English): Overlay networks connect multiple Docker daemons together and enable swarm services
	//to communicate with each other.
	//
	//You can also use overlay networks to facilitate communication between a swarm service and a standalone container, or
	//between two standalone containers on different Docker daemons.
	//
	//This strategy removes the need to do OS-level routing between these containers.
	//
	//See overlay networks. https://docs.docker.com/network/overlay/
	KNetworkDriveOverlay

	//KNetworkDriveMacVLan (English): Macvlan networks allow you to assign a MAC address to a container, making it appear
	//as a physical device on your network.
	//
	//The Docker daemon routes traffic to containers by their MAC addresses. Using the macvlan driver is sometimes the
	//best choice when dealing with legacy applications that expect to be directly connected to the physical network,
	//rather than routed through the Docker host’s network stack.
	//
	//See Macvlan networks. https://docs.docker.com/network/macvlan/
	KNetworkDriveMacVLan

	//KNetworkDriveNone (English): For this container, disable all networking. Usually used in conjunction with a custom
	//network driver.
	//
	//None is not available for swarm services.
	//
	//See disable container networking. https://docs.docker.com/network/none/
	KNetworkDriveNone
)

var networksDrivers = [...]string{
	"bridge",
	"host",
	"overlay",
	"macvlan",
	"none",
}
