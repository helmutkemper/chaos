# Chaos Test

> This is version 3.0, still under development

This project aims to create chaos testing for microservices, allowing to transform a simple golang test into a chaos 
test.

The focus of this project is to allow the chaos test still in the development of the project and try to solve the 
famous problem, on my machine it works!

The test consists of creating all the necessary infrastructure for the project to work on the developer's machine, 
using docker, and after that, pausing or dropping containers, stopping the data flow in the middle of the process.

Imagine making a microservice where three instances must keep data up to date with each other at all times.

The test allows you to create a container for each instance, simulate the data and leave the data flowing, while the 
containers are paused or restarted at random, pausing the transmission of data unexpectedly, allowing you to capture 
failures such as loss of connection or excessive delay in the transmission of data.


## Basic usage

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

## Example

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

> The `main.go` and `go.mod` files must be contained in the root folder of the project

## Simulate network problems

```go
package mongodbproject

import (
	"context"
	"fmt"
	"github.com/helmutkemper/chaos/factory"
	"go.mongodb.org/mongo-driver/bson"
	"runtime/debug"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	factory.NewContainerNetworkProxy(
		"delay",
		27016,
		"delete_mongo_0:27017",
		10, 100,
		0.0,
	)

	go mongoPopulate(t)

	if !primordial.Monitor(2 * time.Minute) {
		t.Fail()
	}
}

func mongoPopulate(t *testing.T) {
	var err error
	var mongoClient *mongo.Client
	var start = time.Now()

	fmt.Printf("conexão\n")

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://0.0.0.0:27016"))
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Connect(context.Background())
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		panic(string(debug.Stack()))
	}

	type Trainer struct {
		Name string
		Age  int
		City string
	}

	var collection *mongo.Collection
	var totalOfInserts int64 = 1000
	for i := int64(0); i != totalOfInserts; i += 1 {
		collection = mongoClient.Database("test").Collection("trainers")
		ash := Trainer{"Ash", 10, "Pallet Town"}
		_, err = collection.InsertOne(context.Background(), ash)
		if err != nil {
			panic(err)
		}
	}

	var total int64
	if total, err = collection.CountDocuments(context.Background(), bson.M{"name": "Ash", "age": 10, "city": "Pallet Town"}); err != nil {
		panic(err)
	}

	if total != totalOfInserts {
		t.Logf("total of inserts must be %v found %v", totalOfInserts, total)
		t.Fail()
	}

	fmt.Printf("fim\n")
	duration := time.Since(start)
	fmt.Printf("Duration: %v\n\n", duration)
}
```