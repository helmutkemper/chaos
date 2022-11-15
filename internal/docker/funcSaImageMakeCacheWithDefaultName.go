package docker

import "time"

// SaImageMakeCache
//
// English:
//
//	Creates a cached image used as a basis for creating new images.
//
//	 Input:
//	   projectPath: path of the project folder
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
//	 * Usar as funções container.MakeDefaultDockerfileForMeWithInstallExtras() ou
//	   container.MakeDefaultDockerfileForMe().
//
//	Segunda opção:
//
//	 * Criar uma pasta contendo o arquivo Dockerfile a ser usado como base para a criação de novas
//	   imagens;
//	 * Criar seu próprio Dockerfile e em vez de usar `FROM golang:1.16-alpine`, usar o nome da cacge,
//	   por exemplo, `FROM cache:latest`;
func SaImageMakeCacheWithDefaultName(projectPath string, expirationDate time.Duration) (err error) {
	return SaImageMakeCache(projectPath, "cache:latest", expirationDate)
}
