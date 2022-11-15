package docker

// SetGitConfigFile
//
// English:
//
//	Defines the contents of the .gitconfig file
//
//	 Input:
//	   value: .gitconfig file contents
//
// Example:
//
//	var err error
//	var usr *user.User
//	var path string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	path = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(path)
//	if err != nil {
//	  panic(err)
//	}
//
//	var container = ContainerBuilder{}
//	container.SetGitConfigFile(string(file))
//
// Português:
//
//	 Define o conteúdo do arquivo .gitconfig
//
//		 Entrada:
//	    value: conteúdo do arquivo .gitconfig
//
// Exemplo:
//
//	var err error
//	var usr *user.User
//	var path string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	path = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(path)
//	if err != nil {
//	  panic(err)
//	}
//
//	var container = ContainerBuilder{}
//	container.SetGitConfigFile(string(file))
func (e *ContainerBuilder) SetGitConfigFile(value string) {
	e.contentGitConfigFile = value
}
