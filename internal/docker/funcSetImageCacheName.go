package docker

// SetImageCacheName
//
// English::
//
//	Defines the name of the cache image
//
//	 Input:
//	   name: Name of the cached image. (Default: "cache:lastest")
//
// Note:
//
//   - See SaImageMakeCache(), SetCacheEnable(), MakeDefaultDockerfileForMe() and
//     MakeDefaultDockerfileForMeWithInstallExtras() functions
//
// Português:
//
//	Define o nome da imagem cache
//
//	 Entrada:
//	   name: Nome da imagem cacge. (Default: "cache:lastest")
//
// Nota:
//
//   - Veja as funções SaImageMakeCache(), SetCacheEnable(), MakeDefaultDockerfileForMe() e
//     MakeDefaultDockerfileForMeWithInstallExtras()
func (e *ContainerBuilder) SetImageCacheName(name string) {
	e.imageCacheName = name
}
