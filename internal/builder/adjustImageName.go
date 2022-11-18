package builder

import "strings"

// AdjustImageName
//
// Test and adjust image name to name+":"+version_string
func (el DockerSystem) AdjustImageName(imageName string) (name string) {

	if strings.Contains(imageName, ":") == false {
		imageName = imageName + ":latest"
		return imageName
	}

	if strings.HasSuffix(imageName, ":") == true {
		imageName = imageName + "latest"
	}

	return imageName
}
