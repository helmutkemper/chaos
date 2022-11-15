package docker

// SetContainerCommandToRunWhenStartingTheContainer
//
// English:
//
//	Command to run when stating the container (style Dockerfile CMD)
//
//	 Input:
//	   values: List of commands. Eg.: []string{"ls", "-l"}
//
// Português:
//
//	Comando a ser executado quando o container inicia (estilo Dockerfile CMD)
//
//	 Entrada:
//	   values: Lista de comandos. Ex.: []string{"ls", "-l"}
func (e *ContainerBuilder) SetContainerCommandToRunWhenStartingTheContainer(values []string) {
	e.containerConfig.Cmd = values
}
