package docker

// GetImageRepoTags
//
// English:
//
//	Get the list of tags of an image repository.
//
//	 Output:
//	   repoTags: list of image repository tags.
//
// Português:
//
//	Obtém a lista de tags de um repositório de imagens.
//
//	 Saída:
//	   repoTags: lista de tags do repositório de imagens.
func (e *ContainerBuilder) GetImageRepoTags() (repoTags []string) {
	return e.imageRepoTags
}
