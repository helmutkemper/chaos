package docker

import "time"

// GetImageExpirationTime
//
// Português:
//
//	Retorna a data de expiração da imagem, ou seja, a data usada como base para impedir que a mesma
//	imagem seja gerada várias vezes em um único teste.
//
//	 Saída:
//	   expiration: Objeto time.Duration com a data de validade da imagem
//
// English:
//
//	Returns the image's expiration date, that is, the date used as a basis to prevent the same image
//	from being generated multiple times in a single test.
//
//	 Output:
//	   expiration: time.Duration object with the image's expiration date
func (e *ContainerBuilder) GetImageExpirationTime() (expiration time.Duration) {
	return e.imageExpirationTime
}
