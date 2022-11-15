package docker

// GetImageParent
//
//	Retorna o nome da imagem base
//
//	 Saída:
//	   parent: nome da imagem base
//
// English:
//
//	Returns the name of the base image
//
//	 Output:
//	   parent: name of the base image
func (e *ContainerBuilder) GetImageParent() (parent string) {
	return e.imageParent
}
