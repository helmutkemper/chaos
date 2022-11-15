package docker

import (
	"time"
)

// SetInspectInterval
//
// English:
//
//	Defines the container's monitoring interval [optional]
//
//	 Input:
//	   value: time interval between container inspection events
//
// Note:
//
//   - This function has a high computational cost and should be used sparingly.
//   - The captured values are presented by GetLastInspect() and GetChannelOnContainerInspect()
//
// Português:
//
//	Define o intervalo de monitoramento do container [opcional]
//
//	 Entrada:
//	   value: intervalo de tempo entre os eventos de inspeção do container
//
// Nota:
//
//   - Esta função tem um custo computacional elevado e deve ser usada com moderação.
//   - Os valores capturados são apresentados por GetLastInspect() e GetChannelOnContainerInspect()
func (e *ContainerBuilder) SetInspectInterval(value time.Duration) {
	e.inspectInterval = value
}
