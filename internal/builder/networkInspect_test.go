package builder

import (
	"errors"
	"github.com/docker/docker/api/types"
)

func ExampleDockerSystem_NetworkInspect() {

	var err error
	var networkId string
	var dockerSys *DockerSystem

	// English: create a new default client. Please, use: err, dockerSys = factoryDocker.NewClient()
	// Português: cria um novo cliente com configurações padrão. Por favor, usr: err, dockerSys = factoryDocker.NewClient()
	dockerSys = &DockerSystem{}
	dockerSys.ContextCreate()
	err = dockerSys.ClientCreate()
	if err != nil {
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// English: create a network named 'network_delete_before_test'
	// Português: cria uma nova rede de nome 'network_delete_before_test'
	networkId, _, err = dockerSys.NetworkCreate(
		"network_delete_before_test",
		KNetworkDriveBridge,
		"local",
		"10.0.0.0/16",
		"10.0.0.1",
	)
	if err != nil {
		panic(err)
	}

	if networkId == "" {
		err = errors.New("network id was not generated")
		panic(err)
	}

	var inspect types.NetworkResource
	inspect, err = dockerSys.NetworkInspect(networkId)
	if err != nil {
		panic(err)
	}

	if inspect.ID != networkId {
		err = errors.New("wrong network id")
		panic(err)
	}

	if inspect.Name != "network_delete_before_test" {
		err = errors.New("wrong network name")
		panic(err)
	}

	if len(inspect.IPAM.Config) == 0 {
		err = errors.New("wrong network config")
		panic(err)
	}

	if inspect.IPAM.Config[0].Gateway != "10.0.0.1" {
		err = errors.New("wrong network name")
		panic(err)
	}

	if inspect.IPAM.Config[0].Subnet != "10.0.0.0/16" {
		err = errors.New("wrong network name")
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// Output:
	//
}
