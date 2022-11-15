package docker

import "context"

// ContextCreate (English): Create a background context
//
// ContextCreate (Português): Cria um contexto background
func (el *DockerSystem) ContextCreate() {
	el.ctx = context.Background()
}
