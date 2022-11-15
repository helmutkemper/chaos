package docker

// GetBuildFolderPath
//
// English:
//
//	Returns the project path used to mount the image
//
//	 Output:
//	   buildPath: string with the project path
//
// Português:
//
//	Retorna o caminho do projeto usado para montar a imagem
//
//	 Saída:
//	   buildPath: string com o caminho do projeto
func (e *ContainerBuilder) GetBuildFolderPath() (buildPath string) {
	return e.buildPath
}
