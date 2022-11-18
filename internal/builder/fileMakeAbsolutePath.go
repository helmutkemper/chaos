package builder

import "path/filepath"

// FileMakeAbsolutePath (English): Make the relative file path absolute
//
//	filePath: string relative file path
//
// FileMakeAbsolutePath (Português): Transforma o caminho relativo de arquivo em absoluto
//
//	filePath: caminho relativo do arquivo
func (el *DockerSystem) FileMakeAbsolutePath(
	filePath string,
) (
	fileAbsolutePath string,
	err error,
) {

	fileAbsolutePath, err = filepath.Abs(filePath)
	return
}
