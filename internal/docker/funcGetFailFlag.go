package docker

// GetFailFlag
//
// English:
//
//	Returns the fail indicator flag set by the AddFailMatchFlag(), AddFailMatchFlagToFileLog(),
//	AddFilterToFail() functions,
//
//
//	 Output:
//	   fail: true if test failed
//
// Português:
//
//	Retorna o flag indicador de falha definido pelas funções AddFailMatchFlag(),
//	AddFailMatchFlagToFileLog(), AddFilterToFail()
//
//	 Saída:
//	   fail: true se o teste tiver falhado
func (e *ContainerBuilder) GetFailFlag() (fail bool) {
	return e.chaos.foundFail
}
