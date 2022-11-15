package docker

// SetImageBuildOptionsSquash
//
// English:
//
//	Squash the resulting image's layers to the parent preserves the original image and creates a new
//	one from the parent with all the changes applied to a single layer
//
//	 Input:
//	   value: true preserve the original image and creates a new one from the parent
//
// Português:
//
//	Usa o conteúdo dos layers da imagem pai para criar uma imagem nova, preservando a imagem pai, e
//	aplica todas as mudanças a um novo layer
//
//	 Entrada:
//	   value: true preserva a imagem original e cria uma nova imagem a partir da imagem pai
func (e *ContainerBuilder) SetImageBuildOptionsSquash(value bool) {
	e.buildOptions.Squash = value
}
