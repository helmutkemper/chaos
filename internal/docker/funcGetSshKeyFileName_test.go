package docker

import (
	"fmt"
	"os/user"
)

func ExampleContainerBuilder_GetSshKeyFileName() {

	var err error
	var userData *user.User

	userData, err = user.Current()
	if err != nil {
		return
	}

	var fileName string
	var cb = ContainerBuilder{}
	fileName, err = cb.GetSshKeyFileName(userData.HomeDir)
	if err != nil {
		return
	}

	fmt.Printf("name: %v", fileName)

	// Output:
	// name: id_ecdsa
}
