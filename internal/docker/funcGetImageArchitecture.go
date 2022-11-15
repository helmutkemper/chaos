package docker

// GetImageArchitecture
//
// English:
//
//	Returns the architecture of the image.
//
//	 Output:
//	   architecture: image architecture
//
// Português:
//
//	Retorna a arquitetura da imagem.
//
//	 Saída:
//	   architecture: arquitetura da imagem
func (e *ContainerBuilder) GetImageArchitecture() (architecture string) {
	return e.imageArchitecture
}
