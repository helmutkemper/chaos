package docker

import "strings"

// AdjustImageName (English): Test and adjust image name to name+":"+version_string
//
// AdjustImageName (Português): Testa e ajusta o nome da imagem para nome+":"+versão_string
func (el DockerSystem) AdjustImageName(
	imageName string,
) string {

	if strings.Contains(imageName, ":") == false {
		imageName = imageName + ":latest"
		return imageName
	}

	if strings.HasSuffix(imageName, ":") == true {
		imageName = imageName + "latest"
	}

	return imageName
}
