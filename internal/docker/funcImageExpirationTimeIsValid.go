package docker

import "time"

// imageExpirationTimeIsValid
//
// English:
//
//	Detects if the image is within the expiration date.
//
//	 Output:
//	   valid: true, if the image is within the expiry date.
//
// Português:
//
//	Detecta se a imagem está dentro do prazo de validade.
//
//	 Saída:
//	   valid: true, se a imagem está dentro do prazo de validade.
func (e *ContainerBuilder) imageExpirationTimeIsValid() (valid bool) {
	if e.imageExpirationTime == 0 {
		return
	}

	var err error
	_, err = e.ImageInspect()
	if err != nil {
		return
	}

	return e.GetImageCreatedTime().Add(e.GetImageExpirationTime()).After(time.Now())
}
