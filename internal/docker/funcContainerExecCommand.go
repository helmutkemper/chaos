package docker

import "github.com/helmutkemper/util"

// ContainerExecCommand
//
// Português:
//
//	Executa comandos dentro do container.
//
//	 Entrada:
//	   commands: lista de comandos. Ex.: []string{"ls", "-l"}
//
//	 Saída:
//	   exitCode: código de saída do comando.
//	   runing: indica se o comando está rodando.
//	   stdOutput: saída padrão do comando.
//	   stdError: saída de erro do comando.
//	   err: objeto de erro padrão.
//
// English:
//
//	Execute commands inside the container.
//
//	 Input:
//	   commands: command list. Eg: []string{"ls", "-l"}
//
//	 Output:
//	   exitCode: command exit code.
//	   runing: indicates whether the command is running.
//	   stdOutput: standard output of the command.
//	   stdError: error output from the command.
//	   err: standard error object.
func (e *ContainerBuilder) ContainerExecCommand(
	commands []string,
) (
	exitCode int,
	runing bool,
	stdOutput []byte,
	stdError []byte,
	err error,
) {

	if e.containerID == "" {
		err = e.getIdByContainerName()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	exitCode, runing, stdOutput, stdError, err = e.dockerSys.ContainerExecCommand(e.containerID, commands)
	return
}
