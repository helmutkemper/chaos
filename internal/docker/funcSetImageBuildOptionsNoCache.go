package docker

// SetImageBuildOptionsNoCache
//
// English:
//
//	Set image build no cache
//
// Português:
//
//	Define a opção `sem cache` para a construção da imagem
func (e *ContainerBuilder) SetImageBuildOptionsNoCache() {
	e.buildOptions.NoCache = true
}
