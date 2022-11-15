package docker

// GetLastLogs
//
// English:
//
//	Returns the standard container output based on the last ticker cycle defined in
//	SetInspectInterval()
//
//	 Output:
//	   logs: container standard output text
//
// Note:
//
//   - The GetChannelOnContainerInspect() function returns the channel triggered by the ticker when
//     the information is ready for use
//
// Português:
//
//	Retorna a saída padrão do container baseado no último ciclo do ticker definido em
//	SetInspectInterval()
//
//	 Saída:
//	   logs: saída padrão do container
//
// Nota:
//
//   - A função GetChannelOnContainerInspect() retorna o canal disparado pelo ticker quando as
//     informações estão prontas para uso
func (e *ContainerBuilder) GetLastLogs() (logs string) {
	return e.logs
}
