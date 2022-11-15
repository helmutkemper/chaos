package docker

// GetImageComment
//
// English:
//
//	Returns the archived comment of the image
//
//	 Output:
//	   comment: image comment
//
// Português:
//
//	Retorna o comentário arquivado na imagem
//
//	 Saída:
//	   comment: comentário da imagem
func (e *ContainerBuilder) GetImageComment() (comment string) {
	return e.imageComment
}
