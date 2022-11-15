package docker

// GetGitCloneToBuild
//
// English:
//
//	Returns the URL of the repository to clone for image transformation
//
// Note:
//
//   - See the SetGitCloneToBuild() function for more details.
//
// Português:
//
//	Retorna a URL do repositório a ser clonado para a transformação em imagem
//
// Nota:
//
//   - Veja a função SetGitCloneToBuild() para mais detalhes.
func (e *ContainerBuilder) GetGitCloneToBuild() (url string) {
	return e.gitData.url
}
