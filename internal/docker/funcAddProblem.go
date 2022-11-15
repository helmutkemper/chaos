package docker

// addProblem
//
// English: Adds a description of a problem to explain the error to the user.
//
//	Input:
//	  problem: problem explanation
//
// Português: Adiciona a descrição de um problema para explica o erro ao usuário.
//
//	Entrada:
//	  problem: descrição do problema
func (e *ContainerBuilder) addProblem(problem string) {
	if e.problem == "" {
		e.problem = problem
	} else {
		e.problem += "\n" + problem
	}
}
