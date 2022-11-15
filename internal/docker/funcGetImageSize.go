package docker

// GetImageSize
//
// English:
//
//	Returns the image size
//
//	 Output:
//	   size: image size
//
// Português:
//
//	Retorna o tamanho da imagem
//
//	 Saída:
//	   size: tamanho da imagem
func (e *ContainerBuilder) GetImageSize() (size int64) {
	return e.imageSize
}
