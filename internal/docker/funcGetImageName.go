package docker

// GetImageName
//
// English:
//
//	Returns the name of the image.
//
//	 Output:
//	   name: Name of the image
//
// Português:
//
//	Retorna o nome da imagem.
//
//	 Saída:
//	   name: Nome da imagem
func (e *ContainerBuilder) GetImageName() (name string) {
	return e.imageName
}
