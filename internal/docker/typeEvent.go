package docker

type Event struct {
	ContainerName string
	Message       string
	Error         bool
	Done          bool
	Fail          bool
	Metadata      map[string]interface{}
}

func (e *Event) clear() {
	e.ContainerName = ""
	e.Message = ""
	e.Done = false
	e.Error = false
	e.Fail = false
	e.Metadata = make(map[string]interface{})
}
