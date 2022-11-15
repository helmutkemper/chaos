package docker

// GetLastInspect
//
// English:
//
//	Returns the container data based on the last ticker cycle defined in SetInspectInterval()
//
//	 Output:
//	   inspect: Container information such as name, ID, volumes, network, etc.
//
// Note:
//
//   - The GetChannelOnContainerInspect() function returns the channel triggered by the ticker when
//     the information is ready for use
//
// Português:
//
//	Retorna os dados do container baseado no último ciclo do ticker definido em SetInspectInterval()
//
//	 Output:
//	   inspect: Informações sobre o container, como nome, ID, volumes, rede, etc.
//
// Nota:
//
//   - A função GetChannelOnContainerInspect() retorna o canal disparado pelo ticker quando as
//     informações estão prontas para uso
func (e *ContainerBuilder) GetLastInspect() (inspect ContainerInspect) {
	return e.inspect
}
