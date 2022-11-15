package docker

// GetChannelOnContainerReady
//
// English:
//
//	Channel fired when the container is ready for use
//
// Note:
//
//   - This channel expects the container to signal that it is ready, but it does not take into
//     account whether the application contained in the container is ready. For this reason, it is
//     recommended to use SetWaitString()
//
// Português:
//
//	Canal disparado quando o container está pronto para uso
//
// Nota:
//
//   - Este canal espera o container sinalizar que está pronto, mas, ele não considera se a aplicação
//     contida no container está pronta. Por isto, é recomendado o uso de SetWaitString()
func (e *ContainerBuilder) GetChannelOnContainerReady() (channel *chan bool) {
	return e.onContainerReady
}
