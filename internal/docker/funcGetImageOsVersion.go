package docker

// GetImageOsVersion
//
// English:
//
//	Returns the operating system version of the image
//
//	 Output:
//	   osVersion: operating system version of the image
//
// Português:
//
//	Retorna a versão do sistema operacional da imagem
//
//	 Saída:
//	   osVersion: versão do sistema operacional da imagem
func (e *ContainerBuilder) GetImageOsVersion() (osVersion string) {
	return e.imageOsVersion
}
