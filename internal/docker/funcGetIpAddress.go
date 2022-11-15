package docker

// GetIPV4Address
//
// English:
//
//	Returns the last IP read from the container
//
// Note:
//
//   - If the container is disconnected or connected to another network after creation, this
//     information may change
//
// Português:
//
//	Retorna o último IP lido do container
//
// Nota:
//
//   - Caso o container seja desconectado ou conectado a outra rede após a criação, esta informação
//     pode mudar
func (e *ContainerBuilder) GetIPV4Address() (IP string) {
	return e.IPV4Address
}
