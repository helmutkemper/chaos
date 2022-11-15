package docker

type ContainerPullProgressDetail struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
