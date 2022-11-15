package docker

// SetContainerAttachStandardStreamsToTty
//
// English:
//
//	Attach standard streams to tty
//
//	 Entrada:
//	   value: true to attach standard streams to tty
//
// Português:
//
//	Anexa a saída padrão do tty
//
//	 Entrada:
//	   value: true para anexar a saída padrão do tty
func (e *ContainerBuilder) SetContainerAttachStandardStreamsToTty(value bool) {
	e.containerConfig.Tty = value
}
