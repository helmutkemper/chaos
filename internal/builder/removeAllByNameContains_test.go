package builder

import (
	"errors"
)

func ExampleDockerSystem_RemoveAllByNameContains() {
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

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	_, err = dockerSys.NetworkFindIdByName("network_delete_before_test")
	if err == nil {
		err = errors.New("network don't deleted")
		panic(err)
	}
	if err.Error() != "network not found" {
		panic(err)
	}

	// Output:
	//
}
