package docker

import "time"

// SetTimeToStartChaosOnChaosScene
//
// English:
//
//	This function sets a timeout before the chaos test starts, when indicator text is encountered in
//	the standard output.
//
//	 Input:
//	   min: minimum waiting time until chaos test starts
//	   max: maximum waiting time until chaos test starts
//
//	Basically, the idea is that you put at some point in the test a text like, chaos can be
//	initialized, in the container's standard output and the time gives a random character to when the
//	chaos starts.
//
// Note:
//
//   - The following functions are used together during chaos testing:
//     [optional] iotmakerdockerbuilder.ConfigChaosScene()
//
//     Mandatory set:
//     ContainerBuilder.EnableChaosScene()
//     ContainerBuilder.SetTimeOnContainerUnpausedStateOnChaosScene()
//     ContainerBuilder.SetTimeToStartChaosOnChaosScene()
//     ContainerBuilder.SetTimeToRestartThisContainerAfterStopEventOnChaosScene()
//     ContainerBuilder.SetTimeOnContainerPausedStateOnChaosScene()
//     ContainerBuilder.SetTimeBeforeStartChaosInThisContainerOnChaosScene()
//     ContainerBuilder.SetSceneNameOnChaosScene()
//     [optional] ContainerBuilder.ContainerSetDisabePauseOnChaosScene()
//     [optional] ContainerBuilder.ContainerSetDisabeStopOnChaosScene()
//
// Português:
//
//	Esta função define um tempo de espera antes do teste de caos começar, quando o texto indicador é
//	incontrado na saída padrão.
//
//	 Entrada:
//	   min: tempo mínimo de espera até o teste de caos começar
//	   max: tempo máximo de espera até o teste de caos começar
//
//	Basicamente, a ideia é que você coloque em algum ponto do teste um texto tipo, caos pode ser
//	inicializado, na saída padrão do container e o tempo dá um caráter aleatório a quando o caos
//	começa.
//
// Nota:
//
//   - As funções a seguir são usadas em conjunto durante o teste de caos:
//     [opcional] iotmakerdockerbuilder.ConfigChaosScene()
//
//     Conjunto obrigatório:
//     ContainerBuilder.EnableChaosScene()
//     ContainerBuilder.SetTimeOnContainerUnpausedStateOnChaosScene()
//     ContainerBuilder.SetTimeToStartChaosOnChaosScene()
//     ContainerBuilder.SetTimeToRestartThisContainerAfterStopEventOnChaosScene()
//     ContainerBuilder.SetTimeOnContainerPausedStateOnChaosScene()
//     ContainerBuilder.SetTimeBeforeStartChaosInThisContainerOnChaosScene()
//     ContainerBuilder.SetSceneNameOnChaosScene()
//     [opcional] ContainerBuilder.ContainerSetDisabePauseOnChaosScene()
//     [opcional] ContainerBuilder.ContainerSetDisabeStopOnChaosScene()
func (e *ContainerBuilder) SetTimeToStartChaosOnChaosScene(min, max time.Duration) {
	e.chaos.minimumTimeToStartChaos = min
	e.chaos.maximumTimeToStartChaos = max
}
