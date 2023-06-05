package complex_chaos

import (
	"bytes"
	"github.com/helmutkemper/chaos/factory"
	"testing"
	"time"
)

func TestComplexChaos(t *testing.T) {

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		Test(t, "./end")

	// Estrutura do banco MongoDB com replica set
	//
	//           +-------------+
	//           |             |
	//           |   arbiter   |
	//           |   MongoDB   |
	//           |             |
	//           +------+------+
	//                  |
	//        +---------+---------+
	//        |                   |
	// +------+------+     +------+------+
	// |             |     |             |
	// |  replica 1  |     |  replica 2  |
	// |   MongoDB   |     |   MongoDB   |
	// |             |     |             |
	// +-------------+     +-------------+
	//
	mongoDocker := factory.NewContainerFromImage(
		"mongo:latest",
	).
		// Impede que o MongoDB aceite conexão externa diretamente;
		// Cada banco aceitará apenas conexão do container "delay" especificado;
		EnvironmentVar([]string{"bindIp:delete_delay_0"}, []string{"bindIp:delete_delay_1"}, []string{"bindIp:delete_delay_2"}).
		Cmd([]string{"mongod", "--replSet", "rs0"}).
		WaitForFlagTimeout("Waiting for connections", 30*time.Second).

		// Cada arquivo receberá um nome único e não serão sobrescritos em um novo teste
		FailFlag("./bug", "Address already in use", "panic:", "bug:").

		// Habilita o processo de caos
		EnableChaos(1, 1, 1).
		Create("mongo", 3).
		Start()

	var stdOutput []byte
	var err error

	_, _, stdOutput, _, err = mongoDocker.Command(2, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 2: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(1, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 1: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.initiate()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_1:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_2:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh --eval \"rs.status()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if !bytes.Contains(stdOutput, []byte("'rs0'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_0:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_1:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_2:27017'")) {
		t.Logf("replica set, setup failed")
		t.FailNow()
	}

	// Nesse ponto do projeto, a replica set de MongoDB foi configurada com dados efêmeros e está em uma rede docker, mas, a replica set, por regra do MongoDB, só aceita conexão via host name, e host name só funciona na rede docker, por isto o teste deve ser feito em container

	// Estrutura de rede em teste de caos/falha completa                                Log container                        | Log events
	//                                                                                  -------------------------------------+--------------------------------------------------
	//                      +-------------+     +-------------+      +-------------+    03/06/2023 19:14:04: inserted 175000 |
	//                      |             |     |             |      |   control   |    03/06/2023 19:14:06: inserted 176000 |
	//                  +-> |    proxy    | --> |   MongoDB   | <-+- |      of     |    03/06/2023 19:14:07: inserted 177000 |
	//                  |   |             |     |             |   |  |    chaos    |    03/06/2023 19:14:08: inserted 178000 |
	//                  |   +-------------+     +-------------+   |  +-------------+    no data                              | 03/06/2023 19:14:09: pause(): delete_mongo_0
	//                  |   delete_delay_0      delete_mongo_0    |                     no data                              |
	//                  |                                         |                     no data                              | 03/06/2023 19:15:16: unpause(): delete_mongo_0
	// +-------------+  |   +-------------+     +-------------+   |                     03/06/2023 19:15:17: inserted 179000 |
	// |             |  |   |             |     |             |   |                     03/06/2023 19:15:18: inserted 180000 |
	// | golang code | -+-> |    proxy    | --> |   MongoDB   | <-+
	// |             |  |   |             |     |             |   |                     See the example log:
	// +-------------+  |   +-------------+     +-------------+   |                     The log shows MongoDB saving a block of a thousand individual inserts once or twice a second;
	//                  |   delete_delay_1      delete_mongo_1    |                     The first failure happened at 19:12:05 (pause(): delete_mongo_2) and lasted until 19:14:09;
	//                  |                                         |                     The number of saved blocks remains the same, even with a stopped secondary replica;
	//                  |   +-------------+     +-------------+   |                     The second failure happened at 19:14:09 (pause(): delete_mongo_0) and lasted until 19:15:16, however delete_mongo_0 is the "arbiter" bank;
	//                  |   |             |     |             |   |                     The log shows the last block being saved at "03/06/2023 19:14:08: inserted 178000" and then jumps to "03/06/2023 19:15:17: inserted 179000";
	//                  +-> |    proxy    | --> |   MongoDB   | <-+                     Therefore, the replica set was stopped until the event "unpause(): delete_mongo_0" at 19:15:16, therefore, the replica set is limited by the arbiter bank.
	//                      |             |     |             |
	//                      +-------------+     +-------------+                         The standard output of the "delete_mongodbClient_0.log" container will be automatically saved in the ".end" folder
	//                      delete_delay_2      delete_mongo_2                          The pause/stop events will be shown in the standard output of go
	//                            |
	//                            |
	//                            ↓
	// |--------------------- NORMAL NETWORK ---------------------|
	//  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\
	// |             ||             ||             ||             |
	//  \___________/  \___________/  \___________/  \___________/
	//
	//
	//  |-------------------------- SIMULATION NETWORK --------------------------------|
	//  /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\
	// |             |-------|             |-------|             |-------|             |
	//  \___________/         \___________/         \___________/         \___________/
	//  |- package -|- delay -|- package -|- delay -|- package -|- delay -|- package -|
	//
	// Creates a container with the ability to interrupt network packets and simulate a network with problems
	factory.NewContainerNetworkProxy(
		"delay",

		// Uma configuração para cada container proxy
		[]factory.ProxyConfig{
			{
				// Porta de entrada do mundo externo
				LocalPort: 27017,
				// Conexão com elemento passivo, nesse caso, o MongoDB
				Destination: "delete_mongo_0:27017",

				// Tempo mínimo e máximo para atraso entre pacotes
				MinDelay: 1,
				MaxDelay: 1000000,
			},
			{
				// Porta de entrada do mundo externo
				LocalPort: 27017,
				// Conexão com elemento passivo, nesse caso, o MongoDB
				Destination: "delete_mongo_1:27017",

				// Tempo mínimo e máximo para atraso entre pacotes
				MinDelay: 1,
				MaxDelay: 100,
			},
			{
				// Porta de entrada do mundo externo
				LocalPort: 27017,
				// Conexão com elemento passivo, nesse caso, o MongoDB
				Destination: "delete_mongo_2:27017",

				// Tempo mínimo e máximo para atraso entre pacotes
				MinDelay: 1,
				MaxDelay: 100,
			},
		},
	)

	// Container com o projeto de teste arquivado em uma pasta local, "./mongodbClient"
	factory.NewContainerFromFolder(
		"folder:latest",
		"./mongodbClient",
	).
		// Passar a conexão por environment var deixa o código mais organizado
		EnvironmentVar(
			[]string{
				"CONNECTION_STRING=mongodb://delete_delay_0:27017,delete_delay_1:27017,delete_delay_2:27017/?replicaSet=rs0",
			},
		).
		// Monta o dockerfile de forma automática
		MakeDockerfile().
		// Espera o container rodar
		WaitForFlagTimeout("container is running", 10*time.Second).
		FailFlag("./bug", "panic:").
		Create("mongodbClient", 1).
		Start()

	if !primordial.Monitor(10 * time.Minute) {
		t.Fail()
	}
}
