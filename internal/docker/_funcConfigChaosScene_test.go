package docker

import (
	"errors"
	"fmt"
	dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
	"github.com/helmutkemper/util"
	"strconv"
	"time"
	"github.com/nats-io/go-nats"
)

func ExampleConfigChaosScene() {
	var err error
	
	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()
	
	// English: Create a chaos scene named nats_chaos and control the number of containers stopped at the same time
	//
	// Português: Cria uma cena de caos de nome nats_chaos e controla a quantidade de containers parados ao mesmo tempo
	ConfigChaosScene("nats_chaos", 2, 2, 2)
	
	// English: Create a docker network controler
	//
	// Português: Cria um controlador de rede do docker
	var netDocker = &dockerNetwork.ContainerBuilderNetwork{}
	err = netDocker.Init()
	if err != nil {
		util.TraceToLog()
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	
	// English: Create a network named nats_network_delete_after_test, subnet 10.0.0.0/16 and gatway 10.0.0.1
	//
	// Português: Cria uma rede de nome nats_network_delete_after_test, subrede 10.0.0.0/16 e gatway 10.0.0.1
	err = netDocker.NetworkCreate("nats_network_delete_after_test", "10.0.0.0/16", "10.0.0.1")
	if err != nil {
		util.TraceToLog()
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	
	// English: Create a docker container named container_delete_nats_after_test_ + i
	//
	// Português: Cria um container do docker de nome container_delete_nats_after_test_ + i
	for i := 0; i != 3; i += 1 {
		go func(i int, err error) {
			err = mountNatsContainer(i, netDocker)
			if err != nil {
				util.TraceToLog()
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
		}(i, err)
	}
	
	// English: Let the test run for two minutes before closing it
	//
	// Português: Deixa o teste rodar por dois minutos antes de o encerrar
	time.Sleep(2 * 60 * time.Second)
	
	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()
	
	// Output:
	//
}

func mountNatsContainer(loop int, network *dockerNetwork.ContainerBuilderNetwork) (err error) {
	var container = ContainerBuilder{}
	
	// English: print the standard output of the container
	//
	// Português: imprime a saída padrão do container
	//
	// [optional/opcional]
	container.SetPrintBuildOnStrOut()
	
	// English: Sets a validity time for the image, preventing the same image from being remade for a period of time.
	// In some tests, the same image is created inside a loop, and adding an expiration date causes the same image to be used without having to redo the same image at each loop iteration.
	//
	// Português: Define uma tempo de validade para a imagem, evitando que a mesma imagem seja refeita durante um período de tempo.
	// Em alguns testes, a mesma imagem é criada dentro de um laço, e adicionar uma data de validade faz a mesma imagem ser usada sem a necessidade de refazer a mesma imagem a cada interação do loop
	//
	// [optional/opcional]
	container.SetImageExpirationTime(5 * time.Minute)
	
	// English: Link this container to a chaos scene for greater control
	//
	// Português: Vincula este container a uma cena de caos para maior controle
	container.SetSceneNameOnChaosScene("nats_chaos")
	
	// English: Set image name for docker pull
	//
	// Português: Define o nome da imagem para o docker pull
	container.SetImageName("nats:latest")
	
	// English: set a container name
	//
	// Português: Define o nome do container
	container.SetContainerName("container_delete_nats_after_test_" + strconv.Itoa(loop))
	
	// English: Links the container to the previously created network
	//
	// Português: Vincula o container a rede criada previamente
	container.SetNetworkDocker(network)
	
	// English: Defines a wait for text, where the text must appear in the container's standard output to proceed [optional]
	//
	// Português: Define uma espera por texto, onde o texto deve aparecer na saída padrão do container para prosseguir [opcional]
	container.SetWaitStringWithTimeout("Listening for route connections on 0.0.0.0:6222", 10*time.Second)
	
	// English: Exposes nats listening port to host
	//
	// Português: Expõe a porta de escuta do nats para o host
	container.AddPortToExpose("4222")
	
	// English: Defines the probability of the container restarting and changing the IP address in the process.
	//
	// Português: Define a probalidade do container reiniciar e mudar o endereço IP no processo.
	container.SetRestartProbability(0.9, 1.0, 5)
	
	// English: Defines a time window used to start chaos testing after container initialized
	//
	// Português: Define uma janela de tempo usada para começar o teste de caos depois do container inicializado
	container.SetTimeToStartChaosOnChaosScene(2*time.Second, 5*time.Second)
	
	// English: Sets a time window used to release container restart after the container has been initialized
	//
	// Português: Define uma janela de tempo usada para liberar o reinício do container depois do container ter sido inicializado
	container.SetTimeBeforeStartChaosInThisContainerOnChaosScene(2*time.Second, 5*time.Second)
	
	// English: Defines a time window used to pause the container
	//
	// Português: Define uma janela de tempo usada para pausar o container
	container.SetTimeOnContainerPausedStateOnChaosScene(2*time.Second, 5*time.Second)
	
	// English: Defines a time window used to unpause the container
	//
	// Português: Define uma janela de tempo usada para remover a pausa do container
	container.SetTimeOnContainerUnpausedStateOnChaosScene(2*time.Second, 5*time.Second)
	
	// English: Sets a time window used to restart the container after stopping
	//
	// Português: Define uma janela de tempo usada para reiniciar o container depois de parado
	container.SetTimeToRestartThisContainerAfterStopEventOnChaosScene(2*time.Second, 5*time.Second)
	
	// English: Enable chaos test
	//
	// Português: Habilita o teste de caos
	container.EnableChaosScene(true)
	
	// English: Initialize the container's control object
	//
	// Português: Inicializa o objeto de controle do container
	err = container.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	// English: Container build and start from image nats:latest. Waits for the text "Listening for route connections on 0.0.0.0:6222" to appear  in the standard container output to proceed
	//
	// Português: Constroe e inicializa o container a partir da imagem nats:latest. Espera o texto "Listening for route connections on 0.0.0.0:6222" aparecer na saída padrão do container para prosseguir
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	var IP string
	IP, err = container.FindCurrentIPV4Address()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	if IP != container.GetIPV4Address() {
		err = errors.New("all ip address must be a samer IP")
		util.TraceToLog()
		panic(err)
	}
	
	// Container "container_delete_nats_after_test_" + loop running and ready for use on this code point on var IP;
	// you can use AddPortToExpose("4222"), to open only ports defineds inside code;
	// you can use AddPortToChange("4222", "1111") to open only ports defineds inside code and change port 4222 to port 1111;
	
	// English: Starts container monitoring at two second intervals. This functionality monitors the container's standard output and generates the log defined by the SetCsvLogPath() function.
	//
	// Português: Inicializa o monitoramento do container com intervalos de dois segundos. Esta funcionalidade monitora a saída padrão do container e gera o log definido pela função SetCsvLogPath().
	// StartMonitor() é usado durante o teste de caos e na geração do log de desempenho do container.
	// [optional/opcional]
	container.StartMonitor()
	
	return
}

func natsPublisher() {
	var err error
	var connection *nats.Conn
	connection, err = nats.Connect("nats://10.0.0.2:4222,nats://10.0.0.3:4222,nats://10.0.0.4:4222")
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	go func() {
		var counter = 0
		for i := 0; i != 100; i += 1 {
			time.Sleep(1 * time.Second)
			counter++
			
			err = connection.Publish("counterChannel", []byte("counter: "+strconv.Itoa(counter)))
			if err != nil {
				util.TraceToLog()
				panic(err)
			}
		}
		
		connection.Close()
	}()
}

func natsReceiver() {
	var tout = time.NewTimer()
	var err error
	var connection *nats.Conn
	connection, err = nats.Connect("nats://10.0.0.2:4222,nats://10.0.0.3:4222,nats://10.0.0.4:4222")
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	go func() {
		_, err = connection.Subscribe("counterChannel", func(m *nats.Msg) {
			fmt.Println(string(m.Data))
		})
		if err != nil {
			util.TraceToLog()
			panic(err)
		}
	}()
}

func testNatsConnection(connString string) (err error) {
	var connection *nats.Conn
	connection, err = nats.Connect("nats://10.0.0.2:4222,nats://10.0.0.3:4222,nats://10.0.0.4:4222")
	
	go func() {
		var counter = 0
		for i := 0; i != 100; i += 1 {
			time.Sleep(1 * time.Second)
			counter++
			
			err = connection.Publish("counterChannel", []byte("counter: "+strconv.Itoa(counter)))
		}
		
		connection.Close()
	}()
	
	// Simple Publisher
	err = connection.Publish("foo", []byte("Hello World"))
	
	// Simple Async Subscriber
	connection.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	
	// Responding to a request message
	connection.Subscribe("request", func(m *nats.Msg) {
		m.Respond([]byte("answer is 42"))
	})
	
	// Simple Sync Subscriber
	sub, err := connection.SubscribeSync("foo")
	m, err := sub.NextMsg(3 * time.Second)
	
	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	sub, err := connection.ChanSubscribe("foo", ch)
	msg := <-ch
	
	// Unsubscribe
	sub.Unsubscribe()
	
	// Drain
	sub.Drain()
	
	// Requests
	msg, err := connection.Request("help", []byte("help me"), 10*time.Millisecond)
	
	// Replies
	connection.Subscribe("help", func(m *nats.Msg) {
		connection.Publish(m.Reply, []byte("I can help!"))
	})
	
	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	connection.Drain()
	
	// Close connection
	connection.Close()
}
