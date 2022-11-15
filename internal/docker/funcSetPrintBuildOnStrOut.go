package docker

// SetPrintBuildOnStrOut
//
// English:
//
//	Prints the standard output used when building the image or the container to the standard output
//	of the log.
//
// Português:
//
//	Imprime a saída padrão usada durante a construção da imagem ou do container no log.
func (e *ContainerBuilder) SetPrintBuildOnStrOut() {
	e.printBuildOutput = true
}
