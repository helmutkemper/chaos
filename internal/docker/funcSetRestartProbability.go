package docker

// SetRestartProbability
//
// English:
//
//	Set the restart probability and the probability of changing the ip of the container when it restarts.
//
//	Input:
//	  restartProbability: Probability of restarting a container during the chaos test.
//	  restartChangeIpProbability: Probability of changing the ip of the container when it restarts.
//	  limit: Limit of how many times the container will restart.
//
// Português:
//
//	Define a probabilidade de reiniciar um container durante o teste de caos e a probabilidade de
//	trocar o ip do container quando ele reiniciar.
//
//	 Entrada:
//	   restartProbability: Probabilidade de reiniciar um container durante o teste de caos.
//	   restartChangeIpProbability: Probabilidade de trocar o ip do container quando ele reiniciar.
//	   limit: Limite de quantas vezes o container vai reiniciar.
func (e *ContainerBuilder) SetRestartProbability(restartProbability, restartChangeIpProbability float64, limit int) {
	e.chaos.restartProbability = restartProbability
	e.chaos.restartChangeIpProbability = restartChangeIpProbability
	e.chaos.restartLimit = limit
}
