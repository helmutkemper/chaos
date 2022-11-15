package docker

// GetImageID
//
// English:
//
//	Returns the image ID.
//
//	 Output:
//	   ID: image ID
//
// Português:
//
//	Retorna o ID da imagem.
//
//	 Saída:
//	   ID: ID da imagem
func (e *ContainerBuilder) GetImageID() (ID string) {
	return e.imageID
}
