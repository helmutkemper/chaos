package docker

// GetSuccessFlag
//
// English:
//
//	Get success flag
//
//	 Output:
//	   success: success flag
//
// Português:
//
//	Retorna a bandeira indicadora de sucesso no teste
//
//	 Saída:
//	   success: bandeira indicadora de sucesso no teste
func (e *ContainerBuilder) GetSuccessFlag() (success bool) {
	return e.chaos.foundSuccess
}
