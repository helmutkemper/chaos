package docker

import (
	"github.com/helmutkemper/util"
	"log"
)

// WaitForTextInContainerLog
//
// English:
//
//	Wait for the text to appear in the container's default output
//
//	 Input:
//	   value: searched text
//
//	 Output:
//	   dockerLogs: container's default output
//	   err: standard error object
//
// Português: Espera pelo texto aparecer na saída padrão do container
//
//	Entrada:
//	  value: texto procurado
//
//	Saída:
//	  dockerLogs: saída padrão do container
//	  err: objeto de erro padrão
func (e *ContainerBuilder) WaitForTextInContainerLog(value string) (dockerLogs string, err error) {
	var logs []byte
	logs, err = e.dockerSys.ContainerLogsWaitText(e.containerID, value, log.Writer())
	if err != nil {
		util.TraceToLog()
	}
	return string(logs), err
}
