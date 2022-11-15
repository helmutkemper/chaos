package docker

// SetMetadata
//
// English:
//
//	Sets a list of user-defined data
//
//	 Input:
//	   metadata: map[string]interface{} with user defined data
//
// Português:
//
//	Define uma lista de dados definidos pelo usuário
//
//	 Entrada:
//	   metadata: map[string]interface{} com dados definidos oelo usuário
func (e *ContainerBuilder) SetMetadata(metadata map[string]interface{}) {
	e.metadata = metadata
}
