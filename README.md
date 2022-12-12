# Chaos Test

### Basic usage

This example turns a git repository into a container

```go
package main

import (
	"github.com/helmutkemper/chaos/factory"
	"log"
	"time"
)

func main() {
	primordial := factory.NewPrimordial()
	primordial.NetworkCreate("mongo", "10.0.0.0/16", "10.0.0.1")

	factory.NewManager().
		ContainerFromGit(
			"public:latest",
			"https://github.com/helmutkemper/chaos.public.example.git",
		).
		Ports("tcp", 3000, 3000).
		Create("public", 1).
		Start()

	if !primordial.Monitor(5 * time.Minute) {
		log.Print("fail!")
	}

	primordial.GarbageCollector()
}

```

```shell
  make build
```

> Requires Docker installed
