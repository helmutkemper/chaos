package docker

// ConfigChaosScene
//
// English:
//
//	Add and configure a test scene prevents all containers in the scene from stopping at the same time
//
//	 Input:
//	   sceneName: unique name for the scene
//	   maxStopedContainers: Maximum number of stopped containers
//	   maxPausedContainers: Maximum number of paused containers
//	   maxTotalPausedAndStoppedContainers: Maximum number of containers stopped and paused at the same
//	     time
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
//	Adiciona e configura uma cena de teste impedindo que todos os container da cena parem ao mesmo
//
// tempo
//
//	Entrada:
//	  sceneName: Nome único para a cena
//	  maxStopedContainers: Quantidade máxima de containers parados
//	  maxPausedContainers: Quantidade máxima de containers pausados
//	  maxTotalPausedAndStoppedContainers: Quantidade máxima de containers parados e pausados ao mesmo
//	    tempo
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
func ConfigChaosScene(
	sceneName string,
	maxStopedContainers,
	maxPausedContainers,
	maxTotalPausedAndStoppedContainers int,
) {
	theater.ConfigScene(sceneName, maxStopedContainers, maxPausedContainers, maxTotalPausedAndStoppedContainers)
}
