package docker

// getProbalityNumber
//
// English:
//
//	Returns a random number greater than zero and less than one
//
//	 Output:
//	   probality: Open point floating point number between 0.0 and 1.0
//
// Português:
//
//	Retorna um número aleatório maior do que zero e menor do que um
//
//	 Saída:
//	   probality: Número de ponto flutuante de ponto aberto entre 0.0 e 1.0
func (e *ContainerBuilder) getProbalityNumber() (probality float64) {
	return 1.0 - e.getRandSeed().Float64()
}
