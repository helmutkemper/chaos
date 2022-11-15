package docker

// SetOpenAllContainersPorts
//
// English:
//
//	Automatically exposes all ports listed in the image used to generate the container
//
// Note:
//
//   - The ports exposed in the creation of the container can be defined by
//     SetOpenAllContainersPorts(), AddPortToChange() and AddPortToExpose();
//   - By default, all doors are closed;
//   - The ImageListExposedPorts() function returns all ports defined in the image to be exposed.
//
// Português:
//
//	Expõe automaticamente todas as portas listadas na imagem usada para gerar o container
//
// Nota:
//
//   - As portas expostas na criação do container pode ser definidas por SetOpenAllContainersPorts(),
//     AddPortToChange() e AddPortToExpose();
//   - Por padrão, todas as portas ficam fechadas;
//   - A função ImageListExposedPorts() retorna todas as portas definidas na imagem para serem expostas.
func (e *ContainerBuilder) SetOpenAllContainersPorts() {
	e.openAllPorts = true
}
