package docker

import (
	"time"
)

// SaImageMakeCache
//
// English:
//
//	Creates a cached image used as a basis for creating new images.
//
//	 Input:
//	   projectPath: path of the project folder
//	   cacheName: name of the cache image
//	   expirationDate: expiration date of the image
//
//	 Output:
//	   err: standard object error
//
// The way to use this function is:
//
//	First option:
//
//	 * Create a folder containing the Dockerfile file to be used as a base for creating new images;
//	 * Enable the use of image cache in your projects with the container.SetCacheEnable(true)
//	   function;
//	 * Define the name of the cache image used in your projects, with the
//	   container.SetImageCacheName() function;
//	 * Use container.MakeDefaultDockerfileForMeWithInstallExtras() or
//	   container.MakeDefaultDockerfileForMe() functions.
//
//	Second option:
//
//	 * Create a folder containing the Dockerfile file to be used as a base for creating new images;
//	 * Create your own Dockerfile and instead of using `FROM golang:1.16-alpine`, use the name of the
//	   cacge, eg `FROM cache:latest`;
//
// Português:
//
//	Cria uma imagem cache usada como base para a criação de novas imagens.
//
//	 Input:
//	   projectPath: caminha da pasta do projeto
//	   cacheName: nome da imagem cache
//	   expirationDate: data de expiração da imagem.
//
//	 Output:
//	   err: standard object error
//
// A forma de usar esta função é:
//
//	Primeira opção:
//
//	 * Criar uma pasta contendo o arquivo Dockerfile a ser usado como base para a criação de novas
//	   imagens;
//	 * Habilitar o uso da imagem cache nos seus projetos com a função container.SetCacheEnable(true);
//	 * Definir o nome da imagem cache usada nos seus projetos, com a função
//	   container.SetImageCacheName();
//	 * Usar as funções container.MakeDefaultDockerfileForMeWithInstallExtras() ou
//	   container.MakeDefaultDockerfileForMe().
//
//	Segunda opção:
//
//	 * Criar uma pasta contendo o arquivo Dockerfile a ser usado como base para a criação de novas
//	   imagens;
//	 * Criar seu próprio Dockerfile e em vez de usar `FROM golang:1.16-alpine`, usar o nome da cacge,
//	   por exemplo, `FROM cache:latest`;
func SaImageMakeCache(projectPath, cacheName string, expirationDate time.Duration) (err error) {

	var container = ContainerBuilder{}

	// English: Sets a validity time for the image, preventing the same image from being remade for a
	// period of time.
	// In some tests, the same image is created inside a loop, and adding an expiration date causes the
	// same image to be used without having to redo the same image at each loop iteration.
	//
	// Português: Define uma tempo de validade para a imagem, evitando que a mesma imagem seja refeita
	// durante um período de tempo.
	// Em alguns testes, a mesma imagem é criada dentro de um laço, e adicionar uma data de validade faz
	// a mesma imagem ser usada sem a necessidade de refazer a mesma imagem a cada interação do loop
	container.SetImageExpirationTime(expirationDate)

	// English: print the standard output of the container
	//
	// Português: imprime a saída padrão do container
	container.SetPrintBuildOnStrOut()

	// English: Set image name for docker pull
	//
	// Português: Define o nome da imagem para o docker pull
	container.SetImageName(cacheName)

	// English: Golang project path to be turned into docker image
	//
	// Português: Caminho do projeto em Golang a ser transformado em imagem docker
	container.SetBuildFolderPath(projectPath)

	// English: Initializes the container manager object.
	//
	// Português: Inicializa o objeto gerenciador de container.
	err = container.Init()
	if err != nil {
		return
	}

	// English: Creates an image from a project folder.
	//
	// Português: Cria uma imagem a partir de uma pasta de projeto.
	_, err = container.ImageBuildFromFolder()
	if err != nil {
		return
	}

	return
}
