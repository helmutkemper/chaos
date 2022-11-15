package docker

import "time"

// SetTimeOnContainerUnpausedStateOnChaosScene
//
// English:
//
//	Defines the minimum and maximum times where the container is kept out of the paused state
//
//	 Input:
//	   min: minimum time out of sleep state
//	   max: maximum time out of sleep state
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
//	Define os tempos mínimos e máximos onde o container é mantido fora do estado de pausa
//
//	 Entrada:
//	   min: tempo mínimo fora do estado de pausa
//	   max: tempo máximo fora do estado de pausa
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
func (e *ContainerBuilder) SetTimeOnContainerUnpausedStateOnChaosScene(min, max time.Duration) {
	e.chaos.minimumTimeToUnpause = min
	e.chaos.maximumTimeToUnpause = max
}
