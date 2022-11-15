// Package iotmakerdockerbuilder
//
// English:
//
// # Golang and Docker in a simple way
//
// # Documentation in progress
//
// This package facilitates the use of docker containers by golang code, enabling the creation of unit tests involving
// containers in linear and chaos scenarios, enabling the development of microservices and failure simulation.
//
// Português: Golang e Docker de forma simples.
//
// Este pacote facilita o uso de containers docker por código golang, possibilitando a criação de testes unitários
// envolvendo containers em cenários linear e de caos possibilitando o desenvolvimento de microsserviços e simulação de
// falha.
//
// # Transforme teste unitário em cenário de caos
//
// A criação de microsserviços requerem uma nova abordagem de testes, onde nem sempre, os testes unitários são fáceis de
// fazer.
//
// Imagine um microsserviço simples, uma simples comunicação gRPC entre duas instâncias do mesmo serviço.
//
// Como fazer um simples teste para saber se eles se conectam?
//
// Este módulo tem a finalidade de resolver este problema, adicionando ao código golang de teste a capacidade de criar
// vários elementos docker de forma muito rápida no meio dos testes unitários.
//
// Imagine poder criar uma rede docker, apontar para uma pasta contendo o projeto e subir quantos containers quiser, com a
// capacidade de gerar relatórios e simular falhas de comunicação aleatórias com algumas poucas linhas de código.
//
// # Criando uma rede docker
//
// A rede é opcional e permite controlar melhor o endereço IP de cada instância do serviço em teste, além de permitir
// isolar a comunicação entre eles.
//
// Exemplo de código para criação de rede
//
//	package code
//	import (
//		dockerBuilder "github.com/helmutkemper/iotmaker.docker.builder"
//		dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
//		"log"
//		"testing"
//	)
//
//	func TestCode(t *testing.T) {
//		var err error
//		var netDocker = dockerNetwork.ContainerBuilderNetwork{}
//		err = netDocker.Init()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//
//		// create a network named cache_delete_after_test, subnet 10.0.0.0/16 e gatway 10.0.0.1
//		err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//	}
//
// Uma vez criada a rede, cada instância do serviço adicionada ao docker ganhará um endereço IP seguindo a ordem de criação
// da instância.
//
// Por exemplo, a primeira instância criada irá para o endereço `10.0.0.2` e a seguinte irá para o endereço `10.0.0.3`, e
// assim por diante.
//
// Uma vez criada a rede, basta usar o comando `container.SetNetworkDocker(&netDocker)` e a mesma será ligada a nova rede
// de forma transparente.
//
// Caso queira trocar o IP de uma instância, para simular uma troca de IP aleatória, basta rodar o comando
// `container.NetworkChangeIp()` e a instância terá o seu IP trocado para o próximo IP da lista.
//
// # Subindo um container baseado em uma imagem pública
//
// Imagine que o seu projeto necessita de um container `nats:latest` para rodar, logo temos o código golang:
//
//	package code
//	import (
//		dockerBuilder "github.com/helmutkemper/iotmaker.docker.builder"
//		dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
//		"log"
//		"testing"
//	)
//
//	func TestCode(t *testing.T) {
//		var err error
//		var netDocker = dockerNetwork.ContainerBuilderNetwork{}
//		err = netDocker.Init()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//
//		// Create a network named cache_delete_after_test, subnet 10.0.0.0/16 e gateway 10.0.0.1
//		// Use the term "delete" to enable the function "dockerBuilder.SaGarbageCollector()", which will search for and remove
//		// all docker elements with the term "delete" contained in the name. For example, network, image, container and
//		// volumes.
//		err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//
//		// Create a container
//		var container = dockerBuilder.ContainerBuilder{}
//		// Set image name for docker pull
//		container.SetImageName("nats:latest")
//		// Expose nats port [optional]
//		container.AddPortToExpose("4222")
//		// Link container and network [optional] (next ip address is 10.0.0.2)
//		container.SetNetworkDocker(&netDocker)
//		// Set a container name.
//		// Use the term "delete" to enable the function "dockerBuilder.SaGarbageCollector()", which will search for and remove
//		// all docker elements with the term "delete" contained in the name. For example, network, image, container and
//		// volumes.
//		container.SetContainerName("container_delete_nats_after_test")
//		// Set a waits for the text to appear in the standard container output to proceed [optional]
//		container.SetWaitStringWithTimeout("Listening for route connections on 0.0.0.0:6222", 10*time.Second)
//		// Initialize the container object
//		err = container.Init()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//		// Image nats:latest pull command
//		err = container.ImagePull()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//		// Container build and start from image nats:latest
//		// Waits for the text "Listening for route connections on 0.0.0.0:6222" to appear  in the standard container
//		// output to proceed
//		err = container.ContainerBuildFromImage()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//	}
//
// Como padrão, todos os parâmetros são adicionados primeiro e em seguida o objeto é inicializado, com o comando
// `container.Init()`.
//
// Como este exemplo usa uma imagem pública, o primeiro comando é o comando `container.ImagePull()`, para que a imagem
// definida em `container.SetImageName("nats:latest")` seja baixada.
//
// Logo em seguida, o comando `container.ContainerBuildFromImage()` gera um container de nome
// `container.SetContainerName("container_delete_nats_after_test")` e deixa o código parado até a saída padrão do container
// exibir o texto [opcional] `container.SetWaitStringWithTimeout("Listening for route connections on 0.0.0.0:6222", 10*time.Second)`.
//
// # Subindo um container baseado em uma pasta local com acesso a repositório privado
//
// Esta configuração permite transformar uma pasta local em uma imagem, de forma simples, mesmo se o projeto necessitar acessar
// um repositório git privado, protegido com chave `ssh`
//
//	package code
//
//	import (
//		dockerBuilder "github.com/helmutkemper/iotmaker.docker.builder"
//		dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
//		"log"
//		"testing"
//	)
//
//	func TestCode(t *testing.T) {
//		var err error
//		var netDocker = dockerNetwork.ContainerBuilderNetwork{}
//		err = netDocker.Init()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//
//		// Create a network named cache_delete_after_test, subnet 10.0.0.0/16 e gateway 10.0.0.1
//		// Use the term "delete" to enable the function "dockerBuilder.SaGarbageCollector()", which will search for and remove
//		// all docker elements with the term "delete" contained in the name. For example, network, image, container and
//		// volumes.
//		err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//
//		// Create a container
//		container = dockerBuilder.ContainerBuilder{}
//		// Adds an expiration date, in order to prevent the creation of the same image multiple times in a row during the
//		// same test [optional]
//		container.SetImageExpirationTime(5 * time.Minute)
//		// Link container and network [optional] (next ip address is 10.0.0.2)
//		container.SetNetworkDocker(netDocker)
//		// Print the container's standard output to golang's standard output
//		container.SetPrintBuildOnStrOut()
//		// Enables the use of the "cache:latest" image [optional].
//		// To prevent an image from downloading the same dependency multiple times for each test, you can create an image
//		// named "cache:latest" and use this image as the basis for the test images.
//		container.SetCacheEnable(true)
//		// Determines the name of the image to be created during the test.
//		// Use the term "delete" to enable the function "dockerBuilder.SaGarbageCollector()", which will search for and remove
//		// all docker elements with the term "delete" contained in the name. For example, network, image, container and
//		// volumes.
//		container.SetImageName("data_rand_pod_image:latest")
//		// Determina o nome do container. Lembre-se que ele deve ser único.
//		// Use the term "delete" to enable the function "dockerBuilder.SaGarbageCollector()", which will search for and remove
//		// all docker elements with the term "delete" contained in the name. For example, network, image, container and
//		// volumes.
//		container.SetContainerName("delete_data_rand_pod_container")
//		// Determines the path of the folder where your project is located.
//		container.SetBuildFolderPath("./project_folder")
//		// Enables the creation of the "Dockerfile-iotmaker" file automatically, as long as the "main.go" and "go.mod" files
//		// are in the project root.
//		container.MakeDefaultDockerfileForMe()
//		// Defines a list of private repositories used in the project. Separate values by a comma.
//		container.SetGitPathPrivateRepository("github.com/helmutkemper")
//		// Copy the "~/.ssh/id_rsa.pub" and "~/.ssh/known_hosts" credentials into the container.
//		// The automatic creation of the container is done in two steps and the credentials are erased when the first image
//		// is erased.
//		err = container.SetPrivateRepositoryAutoConfig()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//		// Set a waits for the text to appear in the standard container output to proceed [optional]
//		container.SetWaitStringWithTimeout("data rand container started", 10*time.Second)
//		// It links a folder/file contained on the computer where the test runs and a folder/file contained in the container
//		// [optional]
//		err = container.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./memory/container", "/containerMemory")
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//		// It links a folder/file contained on the computer where the test runs and a folder/file contained in the container
//		// [optional]
//		err = container.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./memory/config", "/config")
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//		// Initialize the container object
//		err = container.Init()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//		// Generates an image from a project folder
//		_, err = container.ImageBuildFromFolder()
//		if err != nil {
//				log.Printf("error: %v", err)
//				t.Fail()
//		}
//		// Container build and start from image nats:latest
//		// Waits for the text "data rand container started" to appear  in the standard container
//		// output to proceed
//		err = container.ContainerBuildFromImage()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//	}
//
// Os comandos básicos para a criação de imagem são `container.SetBuildFolderPath("./project_folder")`, para definir a
// pasta local, onde o projeto se encontra, e `container.ImageBuildFromFolder()`, encarregado de transformar o conteúdo da
// pasta em imagem.
//
// Caso haja a necessidade de compartilhar conteúdo local com o container, o comando
// `container.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./memory/config", "/config")` fará a ligação entre
// pastas e arquivos no computador local com o container.
//
// # Criando uma imagem de cache
//
// Em muitos casos de teste, criar uma imagem de cache ajuda a baixar menos dependência na hora de criar as imagens e deixa
// o teste mais rápido.
//
// A forma de fazer isto é bem simples, basta criar uma imagem de nome `cache:latest`.
//
//	package code
//
//	import (
//		dockerBuilder "github.com/helmutkemper/iotmaker.docker.builder"
//		"log"
//		"testing"
//	)
//
//	func TestCache(t *testing.T) {
//		var err error
//
//		// Create a container
//		container = dockerBuilder.ContainerBuilder{}
//		// Adds an expiration date, in order to prevent the creation of the same image multiple times in a row during the
//		// same test [optional]
//		container.SetImageExpirationTime(365 * 24 * time.Hour)
//		// Print the container's standard output to golang's standard output
//		container.SetPrintBuildOnStrOut()
//		// Determines the name of the image to be created during the test.
//		// Use the term "delete" to enable the function "dockerBuilder.SaGarbageCollector()", which will search for and remove
//		// all docker elements with the term "delete" contained in the name. For example, network, image, container and
//		// volumes.
//		container.SetImageName("cache:latest")
//		// Determines the path of the folder where your project is located.
//		container.SetBuildFolderPath("./cache_folder")
//		// Initialize the container object
//		err = container.Init()
//		if err != nil {
//			log.Printf("error: %v", err)
//			t.Fail()
//		}
//		// Generates an image from a project folder
//		_, err = container.ImageBuildFromFolder()
//		if err != nil {
//				log.Printf("error: %v", err)
//				t.Fail()
//		}
//	}
//
// A criação da cache é usada em paralelo com os comandos `container.SetCacheEnable(true)` e
// `container.MakeDefaultDockerfileForMe()`, onde eles vão usar como base a imagem `cache:latest` e a imagem de cache será
// criada em cima da imagem `golang:1.17-alpine`.
//
// Caso você não tenha prática em criar imagens, use o exemplo abaixo, onde `RUN go get ...` são as dependências usadas por
// você.
//
//	FROM golang:1.17-alpine as builder
//	RUN mkdir -p /root/.ssh/ && \
//			apk update && \
//			apk add openssh && \
//			apk add --no-cache build-base && \
//			apk add --no-cache alpine-sdk && \
//			rm -rf /var/cache/apk/*
//	ARG CGO_ENABLED=0
//
//	RUN go get ...
//	RUN go get ...
//	RUN go get ...
//	RUN go get ...
//	RUN go get ...
//
// # Usando repositórios privados
//
// Caso seus projetos necessitem usar repositórios privados, o comando `container.MakeDefaultDockerfileForMe()` sempre faz
// a criação da imagem em duas etapas e as credencias de segurança ficam na primeira etapa, descartada ao final do processo,
// evitando uma cópia das suas credencias de segurança em uma imagem pública.
//
// O comando `container.SetPrivateRepositoryAutoConfig()` copia as suas credenciais de segurança padrão `~/.ssh/id_rsa.pub`,
// `~/.ssh/known_hosts` e `~/.gitconfig`
//
// Em seguida, devemos informar os repositórios privados com o comando
// `container.SetGitPathPrivateRepository("github.com/user1,github.com/user2")`.
//
// Caso você tenha problema em baixar repositórios privados, adicione o código abaixo ao arquivo `~/.gitconfig`
//
//	[core]
//					autocrlf = input
//	[url "ssh://git@github.com/"]
//					insteadOf = https://github.com/
//	[url "git://"]
//					insteadOf = https://
//
// Para quem não tem prática em processo de build em duas etapas, na primeira etapa é criada uma imagem grande com todas
// as dependências e programas necessários para o processador construção do código. Porém, ao final do processo, apenas o
// binário gerado na primeira etapa é copiado para uma imagem nova, o que deixa a imagem final pequena.
package docker
