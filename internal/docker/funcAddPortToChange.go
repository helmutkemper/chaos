package docker

import (
	dockerfileGolang "github.com/helmutkemper/iotmaker.docker.builder.golang.dockerfile"
)

// AddPortToChange
//
// English:
//
//	Defines a new port to be exposed on the network and links with the port defined in the image
//
//	 Input:
//	   imagePort: port defined in the image, in the form of a numeric string
//	   newPort: new port value to be exposed on the network
//
// Nota:
//
//   - The ports exposed in the creation of the container can be defined by
//     SetOpenAllContainersPorts(), AddPortToChange() e AddPortToExpose();
//   - By default, all doors are closed;
//   - The ImageListExposedPorts() function returns all ports defined in the image to be exposed.
//
// Português:
//
//	Define uma nova porta a ser exposta na rede e vincula com a porta definida na imagem
//
//	 Entrada:
//	   imagePort: porta definida na imagem, na forma de string numérica
//	   newPort: novo valor da porta a se exposta na rede
//
// Nota:
//
//   - As portas expostas na criação do container pode ser definidas por SetOpenAllContainersPorts(),
//     AddPortToChange() e AddPortToExpose();
//   - Por padrão, todas as portas ficam fechadas;
//   - A função ImageListExposedPorts() retorna todas as portas definidas na imagem para serem
//     expostas.
func (e *ContainerBuilder) AddPortToChange(imagePort string, newPort string) {
	if e.changePorts == nil {
		e.changePorts = make([]dockerfileGolang.ChangePort, 0)
	}

	e.changePorts = append(
		e.changePorts,
		dockerfileGolang.ChangePort{
			OldPort: imagePort,
			NewPort: newPort,
		},
	)
}
