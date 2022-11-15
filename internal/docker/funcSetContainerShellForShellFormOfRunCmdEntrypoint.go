package docker

// SetContainerShellForShellFormOfRunCmdEntrypoint
//
// English:
//
//	shell for shell-form of run cmd entrypoint
//
// Português:
//
//	define o terminal (shell) para executar o entrypoint
func (e *ContainerBuilder) SetContainerShellForShellFormOfRunCmdEntrypoint(values []string) {
	e.containerConfig.Shell = values
}
