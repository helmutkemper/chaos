package docker

import "time"

// SetImageExpirationTime
//
// English:
//
//	Sets image validity time, preventing image build more than once within a certain period of time.
//
//	 Input:
//	   expiration: Image expiration time
//
// Note:
//
//   - This feature prevents creation of the same image when the test uses loops to generate multiple
//     containers from the same image.
//
// Português:
//
//	Define o tempo de validade da imagem, evitando o build da imagem mais de uma vez dentro de um
//	certo período de tempo.
//
//	 Entrada:
//	   expiration: Tempo de validade da imagem
//
// Nota:
//
//   - Esta funcionalidade impede a criação da mesma imagem quando o teste usa laços para gerar vários
//     containers da mesma imagem.
func (e *ContainerBuilder) SetImageExpirationTime(expiration time.Duration) {
	e.imageExpirationTime = expiration
}
