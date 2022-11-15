package docker

// GetImageAuthor
//
// English:
//
//	Returns the author of the image.
//
//	 Output:
//	   author: image author
//
// Português:
//
//	Retorna o autor da imagem.
//
//	 Saída:
//	   author: autor da imagem
func (e *ContainerBuilder) GetImageAuthor() (author string) {
	return e.imageAuthor
}
