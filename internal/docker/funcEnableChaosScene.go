package docker

// EnableChaosScene
//
// English:
//
//	Enables chaos functionality in containers.
//
//	 Input:
//	   enable: enable chaos manager
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
//	Habilita a funcionalidade de caos nos containers.
//
//	 Entrada:
//	   enable: habilita o gerenciador de caos
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
func (e *ContainerBuilder) EnableChaosScene(enable bool) {
	e.chaos.enableChaos = enable
}
