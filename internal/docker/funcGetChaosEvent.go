package docker

// GetChaosEvent
//
// English:
//
//	 Returns channel of Chaos Winds.
//
//	  Output:
//	    eventChannel: chaos event channel
//	      ContainerName: container name
//		     Message: Error message or event suffered by container, as aborted by system.
//		     Error: True if there is an error
//		     Done: Trick for when the chaos roll was successful. See the AddSuccessMatchFlag() function
//		     Fail: True for when the chaos test failed. See the functions AddFailMatchFlag(),
//		           AddFailMatchFlagToFileLog(), AddFilterToFail()
//		     Metadata: Data defined by the SetMetadata() function
//
// Português:
//
//	 Retorna o canal de ventos de caos.
//
//	  Saída:
//	    eventChannel: canal de eventos de caos
//	      ContainerName: Nome do container
//		     Message: Mensagem de erro ou evento sofrido pelo container, como abortado pelo sistema.
//		     Error: True se houver erro
//		     Done: True para quando o teste de caos foi bem sucessido. Veja a função AddSuccessMatchFlag()
//		     Fail: True para quando o teste de caos falhou. Veja as funções AddFailMatchFlag(),
//		           AddFailMatchFlagToFileLog(), AddFilterToFail()
//		     Metadata: Dados definidos pela função SetMetadata()
func (e *ContainerBuilder) GetChaosEvent() (eventChannel chan Event) {
	return e.chaos.event
}
