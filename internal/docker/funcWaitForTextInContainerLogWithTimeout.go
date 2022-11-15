package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"time"
)

// WaitForTextInContainerLogWithTimeout
//
// English:
//
//	Wait for the text to appear in the container's default output
//
//	 Input:
//	   value: searched text
//	   timeout: wait timeout
//
//	 Output:
//	   dockerLogs: container's default output
//	   err: standard error object
//
// Português:
//
//	Espera pelo texto aparecer na saída padrão do container
//
//	 Entrada:
//	   value: texto procurado
//	   timeout: tempo limite de espera
//
//	 Saída:
//	   dockerLogs: saída padrão do container
//	   err: objeto de erro padrão
func (e *ContainerBuilder) WaitForTextInContainerLogWithTimeout(value string, timeout time.Duration) (dockerLogs string, err error) {
	var logs []byte
	logs, err = e.dockerSys.ContainerLogsWaitTextWithTimeout(e.containerID, value, timeout, log.Writer())
	if err != nil {
		util.TraceToLog()
	}
	return string(logs), err
}
