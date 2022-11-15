package docker

// GetImageCacheName
//
// English:
//
//	Returns the name of the image cache used to create the image
//
//	 Output:
//	   name: name of the image cache
//
// Português:
//
//	Devolve o nome da imagem cache usada para criar a imagem
//
//	 Saída:
//	   name: nome da imagem cache
func (e *ContainerBuilder) GetImageCacheName() (name string) {
	return e.imageCacheName
}
