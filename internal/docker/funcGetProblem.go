package docker

// GetProblem
//
// English:
//
//	Return problem description when possible.
//
//	 Output:
//	   problem: descrição do problema
//
// Português:
//
//	Retorna a descrição do problema, quando possível
//
//	 Saída:
//	   problem: descrição do problema
func (e *ContainerBuilder) GetProblem() (problem string) {
	return e.problem
}
