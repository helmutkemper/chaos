package docker

// GetChannelOnContainerInspect
//
// English:
//
//	Channel triggered at each ticker cycle defined in SetInspectInterval()
//
// Português:
//
//	Canal disparado a cada ciclo do ticker definido em SetInspectInterval()
func (e *ContainerBuilder) GetChannelOnContainerInspect() (channel *chan bool) {
	return e.onContainerInspect
}
