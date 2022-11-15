package docker

import (
	"time"
)

// SetWaitStringWithTimeout
//
// English:
//
//	Defines a text to be searched for in the container's default output and forces it to wait for the
//	container to be considered ready-to-use
//
//	 Input:
//	   value: text emitted to default output reporting by an expected event
//	   timeout: maximum waiting time
//
// Português:
//
//	Define um texto a ser procurado na saída padrão do container e força a espera do mesmo para se
//	considerar o container como pronto para uso
//
//	 Entrada:
//	   value: texto emitido na saída padrão informando por um evento esperado
//	   timeout: tempo máximo de espera
func (e *ContainerBuilder) SetWaitStringWithTimeout(value string, timeout time.Duration) {
	e.waitString = value
	e.waitStringTimeout = timeout
}
