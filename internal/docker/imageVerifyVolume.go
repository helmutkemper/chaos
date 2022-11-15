package docker

// verify if exposed volume (folder only) defined by user is exposed
// in image
func (el *DockerSystem) ImageVerifyVolume(
	id,
	path string,
) (
	bool,
	error,
) {

	list, err := el.ImageListExposedVolumes(id)
	if err != nil {
		return false, err
	}

	for _, volume := range list {
		if volume == path {
			return true, nil
		}
	}

	return false, nil
}
