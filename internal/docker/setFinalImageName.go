package docker

// SetFinalImageName
//
// English:
//
//	Set a two stage build final image name.
//
//	 Input:
//	   name: name of final image
//
// Português:
//
//	Define o nome da imagem final para construção de imagem em dois estágios.
//
//	 Entrada:
//	   name: nome da imagem final.
func (e *ContainerBuilder) SetFinalImageName(name string) {
	e.autoDockerfile.SetFinalImageName(name)
}
