package docker

// SetSshKnownHostsFile
//
// English:
//
//	Set a sseh knownhosts file
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
//	path = filepath.Join(usr.HomeDir, ".ssh", "known_hosts")
//	file, err = ioutil.ReadFile(path)
//	if err != nil {
//	  panic(err)
//	}
//
//	var container = ContainerBuilder{}
//	container.SetSshKnownHostsFile(string(file))
//
// Português:
//
//	Define o arquivo knownhosts do ssh
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
//	path = filepath.Join(usr.HomeDir, ".ssh", "known_hosts")
//	file, err = ioutil.ReadFile(path)
//	if err != nil {
//	  panic(err)
//	}
//
//	var container = ContainerBuilder{}
//	container.SetSshKnownHostsFile(string(file))
func (e *ContainerBuilder) SetSshKnownHostsFile(value string) {
	e.contentKnownHostsFile = value
}
