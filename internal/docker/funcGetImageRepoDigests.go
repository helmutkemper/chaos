package docker

// GetImageRepoDigests
//
// English:
//
//	Get image reports
//
//	 Output:
//	   repoDigests: image reports
//
// Português:
//
//	Obtém relatórios simplificados da imagem
//
//	 Saída:
//	   repoDigests: relatórios da imagem
//
// English:
//
//	Get image reports
//
//	 Output:
//	   repoDigests: image reports
func (e *ContainerBuilder) GetImageRepoDigests() (repoDigests []string) {
	return e.imageRepoDigests
}
