package docker

import "time"

// GetImageCreatedTime
//
// English:
//
//	Returns the date of creation of the image.
//
//	 Output:
//	   created: Time.Time object with the date of creation of the image.
//
// Português:
//
//	Retorna a data de criação da imagem.
//
//	 Saída:
//	   created: Objeto time.Time com a data de criação da imagem.
func (e *ContainerBuilder) GetImageCreatedTime() (created time.Time) {
	return e.imageCreated
}
