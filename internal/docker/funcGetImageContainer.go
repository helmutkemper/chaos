package docker

// GetImageContainer
//
// English:
//
//	Returns the name of the image used to create the container
//
//	 Output:
//	   imageName: Name of the image used to create the container
//
// Português:
//
//	Retorna o nome da imagem usada para criar o container
//
//	 Saída:
//	   imageName: Nome da imagem usada para criar o container
func (e *ContainerBuilder) GetImageContainer() (imageName string) {
	return e.imageContainer
}
