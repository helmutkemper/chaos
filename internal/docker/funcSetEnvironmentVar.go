package docker

// SetEnvironmentVar
//
// English: Defines environment variables
//
//	value: slice of string containing one environment variable per key
//
// Português: Define as variáveis de ambiente
//
//	value: slice de string contendo um variável de ambiente por chave
func (e *ContainerBuilder) SetEnvironmentVar(value []string) {
	e.environmentVar = value
}
