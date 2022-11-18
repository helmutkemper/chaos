package iotmakerdocker

type ContainerPullProgress struct {
	Stream                     string                      `json:"stream"`
	Status                     string                      `json:"status"`
	ProgressDetail             ContainerPullProgressDetail `json:"progressDetail"`
	ID                         string                      `json:"id"`
	SysStatus                  ContainerPullStatus         `json:"-"`
	ImageName                  string
	SuccessfullyBuildContainer bool
	SuccessfullyBuildImage     bool
}
