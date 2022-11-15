package docker

// GetImageVirtualSize
//
// English:
//
//	Returns the virtual size of the image
//
//	 Output:
//	   virtualSize: virtual size of the image
//
// Português:
//
//	Retorna o tamanho virtual da imagem
//
//	 Saída:
//	   virtualSize: tamanho virtual da imagem
func (e *ContainerBuilder) GetImageVirtualSize() (virtualSize int64) {
	return e.imageVirtualSize
}
