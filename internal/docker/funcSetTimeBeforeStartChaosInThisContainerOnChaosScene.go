package docker

import "time"

// SetTimeBeforeStartChaosInThisContainerOnChaosScene
//
// English:
//
//	Defines the minimum and maximum waiting times before enabling the restart of containers in a chaos
//	scenario
//
//	The choice of time will be made randomly between the minimum and maximum values
//
//	 Input:
//	   min: minimum waiting time
//	   max: maximum wait time
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
//	Define os tempos mínimo e máximos de espera antes de habilitar o reinício dos containers em um
//	cenário de caos
//
//	A escolha do tempo será feita de forma aleatória entre os valores mínimo e máximo
//
//	 Entrada:
//	   min: tempo de espera mínimo
//	   max: tempo de espera máximo
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
func (e *ContainerBuilder) SetTimeBeforeStartChaosInThisContainerOnChaosScene(min, max time.Duration) {
	e.chaos.minimumTimeBeforeRestart = min
	e.chaos.maximumTimeBeforeRestart = max
}
