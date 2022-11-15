package docker

// SetGitPathPrivateRepository
//
// English:
//
//	Path do private repository defined in "go env -w GOPRIVATE=$GIT_PRIVATE_REPO"
//
//	 Input:
//	   value: Caminho do repositório privado. Eg.: github.com/helmutkemper
//
// Português:
//
//	Caminho do repositório privado definido em "go env -w GOPRIVATE=$GIT_PRIVATE_REPO"
//
//	 Entrada:
//	   value: Caminho do repositório privado. Ex.: github.com/helmutkemper
func (e *ContainerBuilder) SetGitPathPrivateRepository(value string) {
	e.gitPathPrivateRepository = value
}
