package docker

import "strings"

// ContainerGetLasNameElement (English): Eliminates the slash '/' in the name of some
// containers
//
// ContainerGetLasNameElement (Português): Elimina a barra '/' do nome de alguns
// containers
func ContainerGetLasNameElement(
	name string,
) string {

	if strings.HasPrefix(name, "/") == true {
		name = strings.Replace(name, "/", "", 1)
	}

	return name
}
