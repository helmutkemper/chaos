package docker

// StopMonitor
//
// English:
//
//	Disable time.Ticker in order to gather performance information from the container in the form of a
//	CSV log and manage chaos control, if it has been enabled.
//
// Note:
//
//   - This function is used in conjunction with the EnableChaosScene(), SetCsvLogPath() and
//     StopMonitor() functions;
//   - StopMonitor() Must be called at the end of the chaos test.
//
// Português:
//
//	Desabilita o time.Ticker com a finalidade de colher informações de desempenho do container na
//	forma de um log CSV e gerencia o controle de caos, caso o mesmo tenha sido habilitado.
//
// Nota:
//
//   - Esta função é usada em conjunto com as funções EnableChaosScene(), SetCsvLogPath() e
//     StopMonitor();
//   - StopMonitor() Deve ser chamado ao final do teste de caos.
func (e *ContainerBuilder) StopMonitor() (err error) {

	e.chaos.monitorRunning = false

	if len(e.chaos.monitorStop) == 0 {
		e.chaos.monitorStop <- struct{}{}
	}

	return
}
