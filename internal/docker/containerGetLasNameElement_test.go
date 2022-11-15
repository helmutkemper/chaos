package docker

import "fmt"

func ExampleContainerGetLasNameElement() {
	name := "/container_mongo"
	fmt.Printf("%v\n", ContainerGetLasNameElement(name))

	name = "db/container_mongo"
	fmt.Printf("%v\n", ContainerGetLasNameElement(name))

	// Output:
	// container_mongo
	// db/container_mongo
}
