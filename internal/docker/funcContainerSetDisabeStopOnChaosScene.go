package docker

// ContainerSetDisabeStopOnChaosScene
//
// English:
//
//	Set the container stop functionality to be disabled when the chaos scene is running
//
//	 Input:
//	   value: true to disable the container stop functionality
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
//	Define se a funcionalidade de parar o container será desabilitada quando a cena de chaos estiver em execução
//
//	 Entrada:
//	   value: true para desabilitar a funcionalidade de parar o container
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
func (e *ContainerBuilder) ContainerSetDisabeStopOnChaosScene(value bool) {
	e.chaos.disableStopContainer = value
}
