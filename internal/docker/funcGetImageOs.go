package docker

// GetImageOs
//
// English:
//
//	Returns the operating system used to create the image.
//
//	 Output:
//	   os: name of the operating system used to create the image
//
// Português:
//
//	Retorna o sistema operacional usado para criar a imagem.
//
//	 Saída:
//	   os: nome do sistema operacional usado para criar a imagem
func (e *ContainerBuilder) GetImageOs() (os string) {
	return e.imageOs
}
