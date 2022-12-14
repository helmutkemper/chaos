# Chaos Test

### Basic usage

### Using a git server

```go
package githubproject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestLinear(t *testing.T) {

  primordial := factory.NewPrimordial().
    NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromGit(
    "server:latest",
    "https://github.com/helmutkemper/chaos.public.example.git",
  ).
    Ports("tcp", 3000, 3000).
    Create("server", 1).
    Start()

  if !primordial.Monitor(3 * time.Minute) {
    t.Fail()
  }
}
```

### Using a docker image

```go
package mongodbproject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestLinear(t *testing.T) {

  primordial := factory.NewPrimordial().
    NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromImage(
    "mongo:latest",
  ).
    Ports("tcp", 27017, 27017).
    EnvironmentVar([]string{"--bind_ip_all"}).
    Create("mongo", 1).
    Start()

  if !primordial.Monitor(3 * time.Minute) {
    t.Fail()
  }
}
```

### Using a local folder

```go
package localFolderProject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestDevOps_Linear(t *testing.T) {
  primordial := factory.NewPrimordial().
    NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromFolder(
    "folder:latest",
    "./project",
  ).
    MakeDockerfile().
    FailFlag("../bug", "panic:", "bug:", "error").
    Create("folder", 3).
    Start()

  if !primordial.Monitor(15 * time.Minute) {
    t.FailNow()
  }
}
```

### Example

```go
package localFolderProject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestDevOps_Linear(t *testing.T) {
  primordial := factory.NewPrimordial().
    NetworkCreate("chaos_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromImage("nats:latest").
    EnableChaos(2,2,2,0.0).
    FailFlag("./bug", "panic:", "bug:", "error").
    SaveStatistics("./").
    Ports("tcp", 4222, 4222, 4223, 4224).
    Create("nats", 3).
    Start()

  factory.NewContainerFromFolder(
    "folder:latest",
    "./project",
  ).
    MakeDockerfile().
    EnableChaos(2,2,2,0.0).
    FailFlag("./bug", "panic:", "bug:", "error").
    SaveStatistics("./").
    Create("folder", 3).
    Start()

  if !primordial.Monitor(60 * time.Minute) {
    t.FailNow()
  }
}
```

> O arquivo `main.go` e `go.mod` devem está contidos da pasta raiz do projeto