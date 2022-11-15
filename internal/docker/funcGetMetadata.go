package docker

// GetMetadata
//
// English:
//
//	Returns a list of user-defined data
//
//	 Output:
//	   metadata: map[string]interface{} with user defined data
//
// Português:
//
//	Retorna uma lista de dados definida oelo usuário
//
//	 Saída:
//	   metadata: map[string]interface{} com dados definidos oelo usuário
func (e *ContainerBuilder) GetMetadata() (metadata map[string]interface{}) {
	return e.metadata
}
