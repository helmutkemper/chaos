package docker

// SetBuildFolderPath
//
// English:
//
//	Defines the path of the folder to be transformed into an image
//
//	 Input:
//	   value: path of the folder to be transformed into an image
//
// Note:
//
//   - The folder must contain a dockerfile file, but since different uses can have different
//     dockerfiles, the following order will be given when searching for the file:
//     "Dockerfile-iotmaker", "Dockerfile", "dockerfile" in the root folder;
//   - If not found, a recursive search will be done for "Dockerfile" and "dockerfile";
//   - If the project is in golang and the main.go file, containing the package main, is contained in
//     the root folder, with the go.mod file, the MakeDefaultDockerfileForMe() function can be used to
//     use a standard Dockerfile file
//
// Português:
//
//	Define o caminho da pasta a ser transformada em imagem
//
//	 Entrada:
//	   value: caminho da pasta a ser transformada em imagem
//
// Nota:
//
//   - A pasta deve conter um arquivo dockerfile, mas, como diferentes usos podem ter diferentes
//     dockerfiles, será dada a seguinte ordem na busca pelo arquivo: "Dockerfile-iotmaker",
//     "Dockerfile", "dockerfile" na pasta raiz.
//   - Se não houver encontrado, será feita uma busca recursiva por "Dockerfile" e "dockerfile"
//   - Caso o projeto seja em golang e o arquivo main.go, contendo o pacote main, esteja contido na
//     pasta raiz, com o arquivo go.mod, pode ser usada a função MakeDefaultDockerfileForMe() para ser
//     usado um arquivo Dockerfile padrão
func (e *ContainerBuilder) SetBuildFolderPath(value string) {
	e.buildPath = value
}
