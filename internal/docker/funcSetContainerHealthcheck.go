package docker

// SetContainerHealthcheck
//
// English:
//
//	Holds configuration settings for the HEALTHCHECK feature.
//
//	 Input:
//	   value: Ponteiro para HealthConfig
//	     Test: Test is the test to perform to check that the container is healthy.
//	           An empty slice means to inherit the default.
//	           The options are:
//	             {}: inherit healthcheck
//	             {"NONE"}: disable healthcheck
//	             {"CMD", args...}: exec arguments directly
//	             {"CMD-SHELL", command}: run command with system's default shell
//
//	     Interval: Interval is the time to wait between checks (Zero means inherit).
//
//	     Timeout: Timeout is the time to wait before considering the check to have hung (Zero means
//	              inherit).
//
//	     StartPeriod: The start period for the container to initialize before the retries starts to
//	                  count down (Zero means inherit).
//
//	     Retries: Retries is the number of consecutive failures needed to consider a container as
//	              unhealthy (Zero means inherit).
//
// Português:
//
//	Adiciona definições de configuração para o recurso HEALTHCHECK.
//
//	  Entrada:
//	    value: Ponteiro para HealthConfig
//	      Test: Test é o teste a ser executado para testar a saúde do container se não for definido,
//	            herda o teste padrão
//	            As opções são:
//	              {}: herda o teste padrão
//	              {"NONE"}: desabilita o healthcheck
//	              {"CMD", args...}: executa os argumentos diretamente
//	              {"CMD-SHELL", command} : executa o comando com shell padrão do sistema
//
//	      Interval: intervalo entre testes (zero herda o valor padrão).
//
//	      Timeout: intervalo de espera antes de considerar o teste com problemas (zero herda o valor
//	               padrão).
//
//	      StartPeriod: tempo de espera pela incialização do container antes dos testes começarem
//	                   (zero herda o valor padrão).
//
//	      Retries: número de testes consecutivos antes de considerar o teste com problemas (zero
//	      herda o valor padrão).
func (e *ContainerBuilder) SetContainerHealthcheck(value *HealthConfig) {
	e.containerConfig.Healthcheck.Test = value.Test
	e.containerConfig.Healthcheck.Interval = value.Interval
	e.containerConfig.Healthcheck.Timeout = value.Timeout
	e.containerConfig.Healthcheck.StartPeriod = value.StartPeriod
	e.containerConfig.Healthcheck.Retries = value.Retries
}
