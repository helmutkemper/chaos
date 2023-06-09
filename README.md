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

```
// Structure in chaos/failure test                                                   golang code - container log          | Chaos events
//                                                                                   -------------------------------------+--------------------------------------------------
//                      +-------------+      +-------------+      +-------------+    03/06/2023 19:14:05: inserted 175000 |
//                      |             |      |   MongoDB   |      |   control   |    03/06/2023 19:14:06: inserted 176000 |
//                  +-> |    proxy    | ---> |             | <-+- |      of     |    03/06/2023 19:14:07: inserted 177000 |
//                  |   |             |      |   Arbiter   |   |  |    chaos    |    03/06/2023 19:14:08: inserted 178000 |
//                  |   +-------------+      +-------------+   |  +-------------+    no data                              | 03/06/2023 19:14:09: pause(): delete_mongo_0 (obs: replica set arbiter)
//                  |   delete_delay_0       delete_mongo_0    |                     no data                              |
//                  |                                          |                     no data                              | 03/06/2023 19:15:16: unpause(): delete_mongo_0 (obs: replica set arbiter)
// +-------------+  |   +-------------+      +-------------+   |                     03/06/2023 19:15:17: inserted 179000 |
// |             |  |   |             |      |   MongoDB   |   |                     03/06/2023 19:15:18: inserted 180000 |
// | golang code | -+-> |    proxy    | ---> |             | <-+
// |             |  |   |             |      |  replica 0  |   |                     See the example log:
// +-------------+  |   +-------------+      +-------------+   |                     The log shows MongoDB saving a block of a thousand individual inserts once or twice a second;
//                  |   delete_delay_1       delete_mongo_1    |                     The first failure happened at 19:12:05 (pause(): delete_mongo_2) and lasted until 19:14:09;
//                  |                                          |                     The number of saved blocks remains the same, even with a stopped secondary replica;
//                  |   +-------------+      +-------------+   |                     The second failure happened at 19:14:09 (pause(): delete_mongo_0) and lasted until 19:15:16, however delete_mongo_0 is the "arbiter" bank;
//                  |   |             |      |   MongoDB   |   |                     The log shows the last block being saved at "03/06/2023 19:14:08: inserted 178000" and then jumps to "03/06/2023 19:15:17: inserted 179000";
//                  +-> |    proxy    | ---> |             | <-+                     Therefore, the replica set was stopped until the event "unpause(): delete_mongo_0" at 19:15:16, therefore, the replica set is limited by the arbiter bank.
//                      |             |      |  replica 1  |
//                      +-------------+      +-------------+                         The standard output of the "delete_mongodbClient_0.log" container will be automatically saved in the ".end" folder
//                      delete_delay_2       delete_mongo_2                          The pause/stop events will be shown in the standard output of go
//                  ↑                    ↑
//                  |                    |    |---------------------------- SIMULATION NETWORK -----------------------------|
//                  |                    |    /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\
//                  |                    +-> |             |-------|             |-------|             |-------|             |
//                  |                         \___________/         \___________/         \___________/         \___________/
//                  |                         |- package -|- delay -|- package -|- delay -|- package -|- delay -|- package -|
//                  |
//                  |
//                  |   |--------------------- NORMAL NETWORK ---------------------|
//                  |    /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\
//                  +-> |             ||             ||             ||             |
//                       \___________/  \___________/  \___________/  \___________/
```


# Basic usage

### MongoDB image

```go
package simple_project_from_image

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/helmutkemper/chaos/factory"
	"github.com/helmutkemper/chaos/internal/manager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime/debug"
	"testing"
	"time"
)

// Test01SimpleProjectFromImage creates a simple MongoDB database installation
//
// This example will show the basic settings of how to create an image-based container and how to expose a port to the
// world.
func Test01SimpleProjectFromImage(t *testing.T) {

	// Purpose of the example: Create an image-based container and expose port 27017
	//
	// +-------------+             +-------------+
	// |             |             |             |
	// | golang code | -> 27017 -> |   MongoDB   |
	// |             |   (open)    |             |
	// +-------------+             +-------------+
	//   172.17.0.1                   10.0.0.2

	// Creates all the necessary infrastructure for the project to function properly.
	primordial := factory.NewPrimordial().
		// NetworkCreate: Creates a network inside docker, isolating the test.
		//                However, it becomes mandatory if you want to use the host name functionality to connect by
		//                container name
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").

		// Test: [optional] Allows the garbage controller to remove any image, network, volume or container created for the
		//       test, both at the start of the test and at the end of the test
		//       t: test framework pointer
		//       pathToSave: Saves, in the "./end" folder, the standard output of all containers removed at the end of the
		//       test
		//       names: [optional] "mongo:6.0.6" removes the image at the end of the test, cleaning up disk space
		//              As a rule, all elements created by the test contain the word `delete` as an identifier of something
		//              created for the test, however, you can pass names of docker elements created for the test that will
		//              be removed at the end of the test. Beware, this is a Contains(docker.element, name) search function
		Test(t, "./end", "mongo:6.0.6")

	// Container factory based on an existing image
	factory.NewContainerFromImage(
		"mongo:6.0.6",
	).
		// [optional] Determines one or more ports to be exposed on the network
		//            Rule: use one line per port and one port per container.
		//            For example: for three containers with port 27017 exposed on ports 27016, 27017 and 27018, use:
		//            Ports("tcp", 27017, 27016, 27017, 27018).
		//            If you need to expose more than one port, for example port 27018, the port used for secondary
		//            replication, repeat the command Ports("tcp", 27018, ..., ...).
		Ports("tcp", 27017, 27017).

		// [optional] Free connection from any address
		EnvironmentVar([]string{"--bind_ip_all"}).

		// [optional] allows saving MongoDB data in local folder
		// Volumes("/data/db", "./data/db").

		// [optional] allows to rewrite configuration file (one entry per container)
		// Volumes("/data/configdb/mongod.conf", "./data/configdb/mongod.conf").

		// [optional] Wait for some text to appear on the container's standard output before proceeding with the code
		WaitForFlagTimeout("Waiting for connections", 30*time.Second).
		// Determines the name of the container and the number of containers to be created
		Create("mongo", 1).
		// initialize the container
		Start()

	// When the standard output of the container prints the text "Waiting for connections" the code will continue at this
	// point
	// At that moment, in the project directory there will be the following files:
	//   report.mongo:6.0.6.md: Project based on  https://github.com/google/osv-scanner security reporting
	//   stats.delete_mongo.0.csv: Container performance and memory consumption report, based on point-in-time data
	//   captures

	// If you want to control the total test time, create a go routine and let the test run in parallel
	go func(t *testing.T, primordial *manager.Primordial) {
		var err error
		var mongoClient *mongo.Client
		var start = time.Now()

		fmt.Printf("connection\n")

		// Create the MongoDB client
		mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://0.0.0.0:27017"))
		if err != nil {
			panic(string(debug.Stack()))
		}

		// Connect to MongoDB
		err = mongoClient.Connect(context.Background())
		if err != nil {
			panic(string(debug.Stack()))
		}

		// Test the connection
		err = mongoClient.Ping(context.Background(), readpref.Primary())
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			panic(string(debug.Stack()))
		}

		// Create a data structure
		type Trainer struct {
			Name string
			Age  int
			City string
		}

		// Insert the data
		var collection *mongo.Collection
		var totalOfInserts int64 = 1000
		for i := int64(0); i != totalOfInserts; i += 1 {
			collection = mongoClient.Database("test").Collection("trainers")
			ash := Trainer{gofakeit.Name(), gofakeit.Number(14, 99), gofakeit.City()}
			_, err = collection.InsertOne(context.Background(), ash)
			if err != nil {
				panic(err)
			}
		}

		// Test the integrity
		var total int64
		if total, err = collection.CountDocuments(context.Background(), bson.M{}); err != nil {
			panic(err)
		}

		if total != totalOfInserts {
			t.Logf("total of inserts must be %v found %v", totalOfInserts, total)
			t.Fail()
		}

		fmt.Printf("end\n")
		duration := time.Since(start)
		fmt.Printf("Duration: %v\n\n", duration)

		// [optional] if you want to end the test before the specified time
		primordial.Done()
	}(t, primordial)

	// Determine test time
	if !primordial.Monitor(30 * time.Minute) {
		t.Fail()
	}
}
```

### Cloning the Polar project

```go
package githubproject

import (
	"bytes"
	"github.com/helmutkemper/chaos/factory"
	"io/fs"
	"os"
	"testing"
	"time"
)

// Test02GithubProjectLinear Clone de project Polar, a project equivalent to Kafka, written in golang, and create 3
// brokers.
//
// This project shows how to clone a git repository, replace the Dockerfile and use the hard drive of the machine where
// testing takes place to exchange data during local testing
//
// Topics:
//   - Clone a git repository and create a container
//   - Change Dockerfile
//   - Copy a folder and create a container
//   - Create a Dockerfile automatically
//   - Change the Dockerfile default setup
//   - Define docker volumes
//
// If you skipped the previous file, it contains basic usage information from the system
func Test02GithubProjectLinear(t *testing.T) {
	var err error

	// test structure
	//
	// +-------------+     +-------------+     +-------------+     +-------------+
	// |             |     |             |     |             |     |             |
	// |  Producer   |     |   Broker    |     |  Consumer   |     |   Golang    |
	// |             | --> |             | --> |  Process    |     |             |
	// |   Event     |     |    Polar    |     |   Event     |     |    code     |
	// |             |     |             |     |             |     |             |
	// +------+------+     +-------------+     +------+------+     +------+------+
	//        |                                       |                   |
	//        |                                       |                   +--> Saves an end of test indicator file (./data/ignore.end.empty)
	//        |                                       |                        Compares sent data with received data
	//        |                                       |
	//        |                                       +--> Save received data to file (./data/ignore.dataReceived.txt)
	//        |
	//        +--> Save the sent data in a file (./data/ignore.dataSent.txt)
	//

	_ = os.Remove("./data/ignore.dataSent.txt")
	_ = os.Remove("./data/ignore.dataReceived.txt")
	_ = os.Remove("./data/ignore.end.empty")

	// clear data after test
	t.Cleanup(func() {
		factory.NewPrimordial().GarbageCollector()
		_ = os.Remove("./data/ignore.dataSent.txt")
		_ = os.Remove("./data/ignore.dataReceived.txt")
		_ = os.Remove("./data/ignore.end.empty")
		_ = os.RemoveAll("./data")
	})

	// Creates the directory where the data used in the test will be saved
	_ = os.Mkdir("./data", fs.ModePerm)

	// Creates all the necessary infrastructure for the project to function properly.
	primordial := factory.NewPrimordial().
		// NetworkCreate: Creates a network inside docker, isolating the test.
		//                However, it becomes mandatory if you want to use the host name functionality to connect by
		//                container name
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").

		// Test: [optional] Allows the garbage controller to remove any image, network, volume or container created for the
		//       test, both at the start of the test and at the end of the test
		//       t: test framework pointer
		//       pathToSave: Saves, in the "./end" folder, the standard output of all containers removed at the end of the
		//       test
		//       names: [optional] "polar:latest" removes the image at the end of the test, cleaning up disk space
		//              As a rule, all elements created by the test contain the word `delete` as an identifier of something
		//              created for the test, however, you can pass names of docker elements created for the test that will
		//              be removed at the end of the test. Beware, this is a Contains(docker.element, name) search function
		Test(t, "./end", "polar:latest")

	//
	// +-------------+     +-------------+     +-------------+
	// |             |     |             |     |             |
	// |  Broker  0  |     |  Broker  1  |     |  Broker  2  |
	// |    Polar    |     |    Polar    |     |    Polar    |
	// |             |     |             |     |             |
	// +------+------+     +------+------+     +------+------+
	//        ↓                   ↓                   ↓
	// -------+---------+--Docker--Network--+---------+-------
	//                  ↑                   ↑
	//           +------+------+     +------+------+
	//           |             |     |             |
	//           |  Producer   |     |  Consumer   |
	//           |   event     |     |   event     |
	//           |             |     |             |
	//           +------+------+     +------+------+
	//

	// Container factory based on an git server
	factory.NewContainerFromGit(
		"polar:latest",
		"https://github.com/polarstreams/polar.git",
	).

		// Replaces or adds files to the project, in the temporary folder, before the image is created.
		ReplaceBeforeBuild("Dockerfile", "./Dockerfile").

		// [optional] Copies the ssh ~/.ssh/*.* files and the ~/.gitconfig file to the first image when use the function
		// MakeDockerfile()
		// For security and disk space reasons, the first image will be discarded and a second image containing only the
		// code binary will be created.
		//PrivateRepositoryAutoConfig().

		// [optional] Path do private repository defined in "go env -w GOPRIVATE=$GIT_PRIVATE_REPO"
		//GitPathPrivateRepository("github.com/helmutkemper").

		// [optional] Path to private ssh key file
		//GitSshPrivateKeyPath("~/.ssh/id_rsa").
		// Password for the private ssh key
		//GitSshPassword("*************").

		// [optional] Token to private repository
		//GitPrivateToken("*************************************************").

		// [optional] Password to private repository
		//GitPassword("*************").

		// Set up the Polar project according to the manual
		// Thanks to Jorge Bay, creator of the Polar project
		// Just an array repeats all settings for all containers. Multiple arrays, uses one array per container
		EnvironmentVar(
			[]string{
				//"BARCO_DEV_MODE=true",
				"POLAR_SHUTDOWN_DELAY_SECS=0",
				"POLAR_CONSUMER_ADD_DELAY_MS=5000",
				"POLAR_SEGMENT_FLUSH_INTERVAL_MS=500",
				"POLAR_BROKER_NAMES=delete_polar_0,delete_polar_1,delete_polar_2",
				"POLAR_ORDINAL=0",
			},
			[]string{
				//"BARCO_DEV_MODE=true",
				"POLAR_SHUTDOWN_DELAY_SECS=0",
				"POLAR_CONSUMER_ADD_DELAY_MS=5000",
				"POLAR_SEGMENT_FLUSH_INTERVAL_MS=500",
				"POLAR_BROKER_NAMES=delete_polar_0,delete_polar_1,delete_polar_2",
				"POLAR_ORDINAL=1",
			},
			[]string{
				//"BARCO_DEV_MODE=true",
				"POLAR_SHUTDOWN_DELAY_SECS=0",
				"POLAR_CONSUMER_ADD_DELAY_MS=5000",
				"POLAR_SEGMENT_FLUSH_INTERVAL_MS=500",
				"POLAR_BROKER_NAMES=delete_polar_0,delete_polar_1,delete_polar_2",
				"POLAR_ORDINAL=2",
			},
		).

		// [optional] Determines one or more ports to be exposed on the network
		//            Rule: * Use one line per port and one port per container.
		//                  * Use Port(tcp/udp, port on Dockerfile, container_0, container_1, ... container_N)
		//            Eg. Ports("tcp", 9250, 9150, 9250, 9350).
		//                Ports("tcp", 9251, 9151, 9251, 9351).
		//                Ports("tcp", 9252, 9152, 9252, 9352)
		//                                   |-------+-------|
		//Ports("tcp", 9250, 9250).                  3
		//Ports("tcp", 9251, 9251).                  |
		//Ports("tcp", 9252, 9252).                  |
		//                                           |
		// Determines the name of the container and  | the number of containers to be created
		Create("polar", 3). // <---------------------+
		// initialize the container
		// Remember: configure, create and start
		Start()

	// create a polar consuming container
	//
	consumer := factory.NewContainerFromFolder(
		"consumer",
		"./consumer",
	).
		MakeDockerfile().
		//DockerfileBuild("/app", "/app/main", "/app/main.go").
		Volumes("/data", "./data").
		Create("consumer", 1).
		Start()

	// create a polar producer container
	factory.NewContainerFromFolder(
		"producer",
		"./producer",
	).
		// [optional] Generate the Dockerfile automatically
		MakeDockerfile().

		// [optional] Define the command "RUN go build -o /app/main /app/main.go"
		// /app: Defines the name of the directory where the application will be copied in the first image. There are two
		// images, the second only receives the final binary.
		// When the command "copy . ." copies the source code, it goes to the app folder and in this case, the main.go file
		// is in the root of the project, so the build command: "RUN go build -o /app/main /app/main.go"
		//DockerfileBuild("/app", "/app/main", "/app/main.go").

		// Share the folder "./data"
		Volumes("/data", "./data").
		Create("producer", 1).
		Start()

	// define a test timeout
	if !primordial.Monitor(15 * time.Minute) {
		t.FailNow()
	}

	// write this file, indicate test end to producer container
	err = os.WriteFile("./data/ignore.end.empty", nil, fs.ModePerm)
	if err != nil {
		t.Logf("write end simulation error: %v", err)
		t.FailNow()
	}

	t.Logf("end simulation signal sent")
	consumer.WaitStatusNotRunning(2 * time.Minute)

	// test data integrity
	var data []byte
	data, err = os.ReadFile("./data/ignore.dataSent.txt")
	if err != nil {
		t.Logf("read sent data log error: %v", err)
		t.FailNow()
	}
	linesSent := bytes.Split(data, []byte("\n"))

	data, err = os.ReadFile("./data/ignore.dataReceived.txt")
	if err != nil {
		t.Logf("read receiceved data log error: %v", err)
		t.FailNow()
	}
	linesReceived := bytes.Split(data, []byte("\n"))

	for kSent := range linesSent {
		var pass = false

		for kReceived := range linesReceived {
			if bytes.Equal(linesSent[kSent], linesReceived[kReceived]) {
				pass = true
				break
			}
		}

		if !pass {
			t.Logf("%s", linesSent[kSent])
			t.Fail()
		}
	}
}
```

### Network with problems

```go
package network_with_problems

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/helmutkemper/chaos/factory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime/debug"
	"testing"
	"time"
)

// This is a test with network with problems simulation
// In case you skipped the previous explanation, it contains the basic knowledge of using the system. More information is added here.
//
// This example will show how to create a container with the ability to disturb the network connection
func TestLinearNetworkWithProblems(t *testing.T) {

	//                                        bindIp:delete_delay_0
	// +-------------+      +-------------+      +-------------+
	// |             |      |             |      |             |
	// | golang code | ---> |    proxy    | ---> |   MongoDB   |
	// |             |      |             |      |             |
	// +-------------+      +-------------+      +-------------+
	//                      delete_delay_1       delete_mongo_1

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		Test(t, "./end")

	factory.NewContainerFromImage(
		"mongo:6.0.6",
	).
		// Limit connection source to MongoDB
		// The network simulation container has the name "delay", the container will be created with name, and the host name, "delete_delay_0"
		EnvironmentVar([]string{"bindIp:delete_delay_0"}).
		Create("mongo", 1).
		Start()

	//
	// |--------------------- NORMAL NETWORK ---------------------|
	//  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\
	// |             ||             ||             ||             |
	//  \___________/  \___________/  \___________/  \___________/
	//
	//
	//  |-------------------------- SIMULATION NETWORK --------------------------------|
	//  /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\
	// |             |-------|             |-------|             |-------|             |
	//  \___________/         \___________/         \___________/         \___________/
	//  |- package -|- delay -|- package -|- delay -|- package -|- delay -|- package -|
	//
	// Creates a container with the ability to interrupt network packets and simulate a network with problems
	factory.NewContainerNetworkProxy(
		"delay",

		// One configuration for each proxy container
		[]factory.ProxyConfig{
			{
				// Port to the outside world
				LocalPort: 27016,
				// Connection with the passive element, in this case, MongoDB
				Destination: "delete_mongo_0:27017",

				// Minimum and maximum time for delay between packets
				// total test time: ~1.8s
				//MinDelay: 0,
				//MaxDelay: 0,

				// Minimum and maximum time for delay between packets
				// total test time: ~1m58
				MinDelay: 100,
				MaxDelay: 130,

				// Minimum and maximum time for delay between packets
				// error: panic: connection(0.0.0.0:27016[-5]) socket was unexpectedly closed: EOF
				//MinDelay: 100,
				//MaxDelay: 140,
			},
		},
	)

	go func() {
		var err error
		var mongoClient *mongo.Client
		var start = time.Now()

		fmt.Printf("connection\n")

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
			ash := Trainer{gofakeit.Name(), gofakeit.Number(14, 99), gofakeit.City()}
			_, err = collection.InsertOne(context.Background(), ash)
			if err != nil {
				panic(err)
			}
		}

		var total int64
		if total, err = collection.CountDocuments(context.Background(), bson.M{}); err != nil {
			panic(err)
		}

		if total != totalOfInserts {
			t.Logf("total of inserts must be %v found %v", totalOfInserts, total)
			t.Fail()
		}

		fmt.Printf("end\n")
		duration := time.Since(start)
		fmt.Printf("Duration: %v\n\n", duration)

		primordial.Done()
	}()

	if !primordial.Monitor(5 * time.Minute) {
		t.Fail()
	}
}
```

### Complex test

```go
package complex_linear

import (
	"bytes"
	"github.com/helmutkemper/chaos/factory"
	"testing"
	"time"
)

// This is a test with replica set creation for MongoDB.
// In case you skipped the previous explanation, it contains the basic knowledge of using the system. More information
// is added here.
//
// This example will show the initial command configurations of the container, terminal access, terminal response
// handling and docker network concepts
func TestComplexLinear(t *testing.T) {

	primordial := factory.NewPrimordial().
		// According to the MongoDB manual, the replica set will only work if the hostname is defined, not being able to
		// use an IP address
		// For the hostname to work correctly, within docker, a network must be created, so do not comment on the network
		// creation
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		// If you want to keep using the "mongo:6.0.6" image, just don't put its name here
		Test(t, "./end", "mongo:6.0.6")

	// MongoDB database structure with replica set
	//
	//           +-------------+
	//           |             |
	//           |   arbiter   |
	//           |   MongoDB   |
	//           |             |
	//           +------+------+
	//                  |
	//        +---------+---------+
	//        |                   |
	// +------+------+     +------+------+
	// |             |     |             |
	// |  replica 1  |     |  replica 2  |
	// |   MongoDB   |     |   MongoDB   |
	// |             |     |             |
	// +-------------+     +-------------+
	//
	mongoDocker := factory.NewContainerFromImage(
		"mongo:6.0.6",
	).
		// The Create() function tells you to create 3 containers, so the first container will have port 27017 directed to
		// port 27016 on the network and so on.
		// If only one port is passed, only the first container will have its port exposed on the network
		// If the container has multiple ports, repeat one line per port
		Ports("tcp", 27017, 27016, 27017, 27018).

		// Imagine that you need to use port 27018, the port used for secondary replication. Repeat the command with the
		// next port.
		// Ports("tcp", 27018, 27019, 27020, 27021).

		// When working with MongoDB, be careful when specifying the IP to be released, because bindIp: 0.0.0.0 will release
		// the connection through the IP 0.0.0.0, while the flag "--bind_ip_all" will release for all IPs, which is very
		// similar, however, IT IS NOT THE SAME THING.
		// bindIp: 0.0.0.0 only accepts connection specified as IP 0.0.0.0, --bind_ip_all accepts IP specified as 127.0.0.1
		// for example.
		EnvironmentVar([]string{"--bind_ip_all"}).

		// Imagine that only container 0 will receive external access
		// The logic is simple: if only one value is passed, it works for all containers, if more than one value is passed,
		// key 0 goes to container 0, key 1 goes to container 1 and so on.
		// EnvironmentVar([]string{"--bind_ip_all"}, []string{}, []string{}).

		// When the container starts, MongoDB needs to receive the "--replSet NAME_REPLICA_SET" flag, as in the
		// example below:
		// $ docker run -p 27017:27017 --name mongo --net mongo_network mongo mongod --replSet rs0
		// In the case of MongoDB, some tutorials omit the name of the shell that will receive the "--replSet rs0" flag,
		// like this:
		// $ docker run -p 27017:27017 --name mongo --net mongo_network mongo --replSet rs0
		// However, in our case, it must be passed, as in the command below
		Cmd([]string{"mongod", "--replSet", "rs0"}).

		// If you need to wait for a success indicator flag, add a text and the system will be stopped waiting for it,
		// however, be careful, and add the text thinking about case sensitive text.
		// This function uses strings.Contains(container.stdOutput, "Waiting for connections") to search for text, so be
		// careful with very short text
		WaitForFlagTimeout("Waiting for connections", 30*time.Second).

		// [optional] Looks for failure indicator texts in the container's standard output and saves the container's
		// standard output in the indicated folder for later analysis. (caution: the text is case sensitive)
		// If you want to do a quick test, use the word "Waiting" and look in the "bug" folder when the container starts
		// running
		// Beware, it looks for string use the function strings.Contains(container.stdOutput, text), because of this the
		// "fail" flag might find the words "maxFailedInitialSyncAttempts" or "failed". For this reason, I put the colon (:)
		FailFlag("./bug", "Address already in use", "panic:", "bug:").

		//Volumes("/etc/mongod.conf", "./conf/mongod_0.conf", "./conf/mongod_1.conf", "./conf/mongod_2.conf").

		// Determines the creation of 3 containers at addresses 10.0.0.2:27017 (exposed port 27016), 10.0.0.3:27017
		// (exposed port 27017), 10.0.0.4:27017 (exposed port 27018), host names delete_mongo_0, delete_mongo_1 and
		// delete_mongo_2
		// If you need to change the host name, use the HostName() function and specify a name for each container. Remember,
		// hostname requires a network attached to the container
		Create("mongo", 3).

		// Although it goes without saying, it's good to remember that the Create() and Start() functions must be the last
		// two functions called
		Start()

	// At this point in the code, the banks are ready to use, however creating replicas requires the use of terminal
	// commands
	// To do this, specify the container key, 0 for the first container created, the command interpreter to be used and
	// the flag indicating that these commands will be sent via text, "-c", that is:
	// `/bin/bash -c "echo Hello World!"`

	// Write terminal commands to turn mongodb into replicaset
	var stdOutput []byte
	var err error

	// To transform the MongoDB of key container 2, delete_mongo_2, into a replica set secondary, it is necessary to
	// access the container through the terminal, activate the MongoDB terminal and pass the command "rs.secondaryOk()"
	// Explanation:
	//   * "/bin/bash": is the linux command interpreter
	//   * "-c": the command will arrive via text string, example: `bin bash -c "echo Hello World!"`
	//   * "mongosh": is the MongoDB command interpreter
	//   * "127.0.0.1:27017": is the connection address on the network. At this point, notice, the command is accessing
	//     the container directly, and inside the container, the port is 27017 and the address is localhost. Do not
	//     confuse internal access, directly in the container with external access.
	//   * "--eval \"rs.secondaryOk_()\"": eval lets you run a javascript command via text, and since it's text within
	//     text, the quotes are escaped.

	// When this happens, the command will return a text containing an error or success indicator in case of missing
	// indicator. The indicators are:
	//   * DeprecationWarning: In MongoDB 6.0.6 can be ignored
	//   * MongoNetworkError: Database connection failure
	//   * TypeError: syntax error

	_, _, stdOutput, _, err = mongoDocker.Command(2, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 2: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// Repeat the same process for key container 1, delete_mongo_1
	// Note: Secondary containers must receive the command "rs.secondaryOk()" before the main container receives the
	// command "rs.initiate()"
	_, _, stdOutput, _, err = mongoDocker.Command(1, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 1: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// Initializes the container key 0, delete_mongo_0, as the replica set arbiter MongoDB instance
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.initiate()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// The lack of an error indicator in itself is a success indicator, but the command will return a json in the
	// following format:
	// {
	//   info2: 'no configuration specified. Using a default configuration for the set',
	//   me: 'delete_mongo_0:27017',
	//   ok: 1
	// }
	//
	// In case you need to process the json you can use regular expression,
	// https://regex101.com/library/sjOfeq?orderBy=MOST_POINTS&page=3&search=json

	// Adds the MongoDB contained in the delete_mongo_1 container as a member of the replica set
	// Notes:
	//   * Since the command is passed via text within text, beware of escaped quotes;
	//   * Inside the docker network, all MongoDB are on port 27017, ports 27016, 27017 and 27018 are the ports exposed
	//     to the world, not in the docker network;
	//   * The host name "delete_mongo_x" only works inside the docker network
	//   * MongoDB does not accept replica set configuration by IP, only by host name
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_1:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// Out of pure curiosity, in case of success, MongoDB returns the json:
	// {
	//   ok: 1,
	//   '$clusterTime': {
	//     clusterTime: Timestamp({ t: 1685399709, i: 1 }),
	//     signature: {
	//       hash: Binary(Buffer.from("0000000000000000000000000000000000000000", "hex"), 0),
	//       keyId: Long("0")
	//     }
	//   },
	//   operationTime: Timestamp({ t: 1685399709, i: 1 })
	// }

	// Add the next MongoDB instance to the replica set
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_2:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// If you want to do one last check, pass the "rs.status()" command to any mongo instance, it should return a json
	// counting "set: 'rs0'" and "name: 'delete_mongo_x:27017'" for each MongoDB instance
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh --eval \"rs.status()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if !bytes.Contains(stdOutput, []byte("'rs0'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_0:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_1:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_2:27017'")) {
		t.Logf("replica set, setup failed")
		t.FailNow()
	}

	// At this point in the project, the MongoDB replica set has been configured with ephemeral data and is on a docker
	// network, with ports 27016, 27017 and 27018 exposed to the world, but the replica set, by MongoDB rule, only
	// accepts connections via host name, and host name only works on the docker network, so the test must be done
	// in a container

	// Create a container from a local folder
	factory.NewContainerFromFolder(
		"folder:latest",
		"./mongodbClient",
	).

		// Automatically mounts the Dockerfile if the "main.go" file is in the root of the project and the "go.mod" file
		// exists, even if it is blank.
		// You can specify the Dockerfile path, if it is not in the root of the project with the
		// DockerfilePath("./path/inside/container/Dockerfile") command
		MakeDockerfile().
		WaitForFlagTimeout("container is running", 10*time.Second).
		FailFlag("./bug", "panic:").
		Create("mongodbClient", 1).
		Start()

	// Let the project run for 5 minutes
	if !primordial.Monitor(5 * time.Minute) {
		t.Fail()
	}
}
```

### Chaos test

```go
package complex_chaos

import (
	"bytes"
	"github.com/helmutkemper/chaos/factory"
	"testing"
	"time"
)

// This example tests the behavior of the MongoDB replica set when an instance fails in production
// In case you skipped the previous explanation, it contains the basic knowledge of using the system. More information
// is added here.
//
// This example shows how to use random crashes on instances and simulate network issues
func TestComplexChaos(t *testing.T) {

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		Test(t, "./end")

	// MongoDB database structure with replica set
	//
	//           +-------------+
	//           |             |
	//           |   arbiter   |
	//           |   MongoDB   |
	//           |             |
	//           +------+------+
	//                  |
	//        +---------+---------+
	//        |                   |
	// +------+------+     +------+------+
	// |             |     |             |
	// |  replica 1  |     |  replica 2  |
	// |   MongoDB   |     |   MongoDB   |
	// |             |     |             |
	// +-------------+     +-------------+
	//
	mongoDocker := factory.NewContainerFromImage(
		"mongo:latest",
	).
		// Prevents MongoDB from accepting external connection directly;
		// Each bank will only accept connections from the specified "delete_delay_x" container;
		EnvironmentVar([]string{"bindIp:delete_delay_0"}, []string{"bindIp:delete_delay_1"}, []string{"bindIp:delete_delay_2"}).
		Cmd([]string{"mongod", "--replSet", "rs0"}).
		WaitForFlagTimeout("Waiting for connections", 30*time.Second).

		// Save container standard output on failure
		// Each file will be given a unique name and will not be overwritten in a new test
		FailFlag("./bug", "Address already in use", "panic:", "bug:").

		// Enables the chaos process
		// Maximum number of stopped containers: 1
		// Maximum number of paused containers: 1
		// Maximum number of stopped and paused containers: 1
		EnableChaos(1, 1, 1).
		Create("mongo", 3).
		Start()

	var stdOutput []byte
	var err error

	_, _, stdOutput, _, err = mongoDocker.Command(2, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 2: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(1, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 1: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.initiate()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_1:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_2:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh --eval \"rs.status()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if !bytes.Contains(stdOutput, []byte("'rs0'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_0:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_1:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_2:27017'")) {
		t.Logf("replica set, setup failed")
		t.FailNow()
	}

	// At this point in the project, the MongoDB replica set has been configured with ephemeral data and is on a docker
	// network, but the replica set, by MongoDB rule, only accepts connection via host name, and host name only works on
	// the docker network, for this test must be done in container

	// Structure in chaos/failure test                                                   Log container                        | Chaos events
	//                                                                                   -------------------------------------+--------------------------------------------------
	//                      +-------------+      +-------------+      +-------------+    03/06/2023 19:14:05: inserted 175000 |
	//                      |             |      |   MongoDB   |      |   control   |    03/06/2023 19:14:06: inserted 176000 |
	//                  +-> |    proxy    | ---> |             | <-+- |      of     |    03/06/2023 19:14:07: inserted 177000 |
	//                  |   |             |      |   Arbiter   |   |  |    chaos    |    03/06/2023 19:14:08: inserted 178000 |
	//                  |   +-------------+      +-------------+   |  +-------------+    no data                              | 03/06/2023 19:14:09: pause(): delete_mongo_0 (obs: replica set arbiter)
	//                  |   delete_delay_0       delete_mongo_0    |                     no data                              |
	//                  |                                          |                     no data                              | 03/06/2023 19:15:16: unpause(): delete_mongo_0 (obs: replica set arbiter)
	// +-------------+  |   +-------------+      +-------------+   |                     03/06/2023 19:15:17: inserted 179000 |
	// |             |  |   |             |      |   MongoDB   |   |                     03/06/2023 19:15:18: inserted 180000 |
	// | golang code | -+-> |    proxy    | ---> |             | <-+
	// |             |  |   |             |      |  replica 0  |   |                     See the example log:
	// +-------------+  |   +-------------+      +-------------+   |                     The log shows MongoDB saving a block of a thousand individual inserts once or twice a second;
	//                  |   delete_delay_1       delete_mongo_1    |                     The first failure happened at 19:12:05 (pause(): delete_mongo_2) and lasted until 19:14:09;
	//                  |                                          |                     The number of saved blocks remains the same, even with a stopped secondary replica;
	//                  |   +-------------+      +-------------+   |                     The second failure happened at 19:14:09 (pause(): delete_mongo_0) and lasted until 19:15:16, however delete_mongo_0 is the "arbiter" bank;
	//                  |   |             |      |   MongoDB   |   |                     The log shows the last block being saved at "03/06/2023 19:14:08: inserted 178000" and then jumps to "03/06/2023 19:15:17: inserted 179000";
	//                  +-> |    proxy    | ---> |             | <-+                     Therefore, the replica set was stopped until the event "unpause(): delete_mongo_0" at 19:15:16, therefore, the replica set is limited by the arbiter bank.
	//                      |             |      |  replica 1  |
	//                      +-------------+      +-------------+                         The standard output of the "delete_mongodbClient_0.log" container will be automatically saved in the ".end" folder
	//                      delete_delay_2       delete_mongo_2                          The pause/stop events will be shown in the standard output of go
	//                  ↑                    ↑
	//                  |                    |    |---------------------------- SIMULATION NETWORK -----------------------------|
	//                  |                    |    /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\         /¯¯¯¯¯¯¯¯¯¯¯\
	//                  |                    +-> |             |-------|             |-------|             |-------|             |
	//                  |                         \___________/         \___________/         \___________/         \___________/
	//                  |                         |- package -|- delay -|- package -|- delay -|- package -|- delay -|- package -|
	//                  |
	//                  |
	//                  |   |--------------------- NORMAL NETWORK ---------------------|
	//                  |    /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\  /¯¯¯¯¯¯¯¯¯¯¯\
	//                  +-> |             ||             ||             ||             |
	//                       \___________/  \___________/  \___________/  \___________/
	//
	// Creates a container with the ability to interrupt network packets and simulate a network with problems
	factory.NewContainerNetworkProxy(
		"delay",

		// One configuration for each proxy container
		[]factory.ProxyConfig{
			{
				// Gateway to the outside world
				LocalPort: 27017,
				// Connection with passive element, in this case MongoDB
				Destination: "delete_mongo_0:27017",

				// Minimum and maximum time for delay between packets
				MinDelay: 1,
				MaxDelay: 1000000,
			},
			{
				// Gateway to the outside world
				LocalPort: 27017,
				// Connection with passive element, in this case MongoDB
				Destination: "delete_mongo_1:27017",

				// Minimum and maximum time for delay between packets
				MinDelay: 1,
				MaxDelay: 100,
			},
			{
				// Gateway to the outside world
				LocalPort: 27017,
				// Connection with passive element, in this case MongoDB
				Destination: "delete_mongo_2:27017",

				// Minimum and maximum time for delay between packets
				MinDelay: 1,
				MaxDelay: 100,
			},
		},
	)

	// Container with test project archived in a local folder, "./mongodbClient"
	factory.NewContainerFromFolder(
		"folder:latest",
		"./mongodbClient",
	).
		// Passing the connection through environment var makes the code more organized
		EnvironmentVar(
			[]string{
				"CONNECTION_STRING=mongodb://delete_delay_0:27017,delete_delay_1:27017,delete_delay_2:27017/?replicaSet=rs0",
			},
		).
		// Mount the dockerfile automatically
		MakeDockerfile().
		// Wait for the container to run
		WaitForFlagTimeout("container is running", 10*time.Second).
		FailFlag("./bug", "panic:").
		Create("mongodbClient", 1).
		Start()

	if !primordial.Monitor(10 * time.Minute) {
		t.Fail()
	}
}
```

# Vulnerability Report Example

> Automatic function

This report is based on an open database and shows known vulnerabilities. Validity: Thu Dec 22 18:29:28 2022

## Path

Path: /scan/go.mod
Type: lockfile

### Packages

| Ecosystem | Package          | Version                           |
|-----------|------------------|-----------------------------------|
| Go        | golang.org/x/net | 0.0.0-20220225172249-27dd8689420f |

### Details:

HTTP/2 server connections can hang forever waiting for a clean shutdown that was preempted by a fatal error. This condition can be exploited by a malicious client to cause a denial of service.

### Affected:

| Ecosystem | Package          |
|-----------|------------------|
| Go        | stdlib           |
| Go        | golang.org/x/net |

| type   | URL                                                                                                                    |
|--------|------------------------------------------------------------------------------------------------------------------------|
| WEB    | [https://groups.google.com/g/golang-announce/c/x49AQzIVX-s](https://groups.google.com/g/golang-announce/c/x49AQzIVX-s) |
| REPORT | [https://go.dev/issue/54658](https://go.dev/issue/54658)                                                               |
| FIX    | [https://go.dev/cl/428735](https://go.dev/cl/428735)                                                                   |

### Details:

An attacker can cause excessive memory growth in a Go server accepting HTTP/2 requests.

HTTP/2 server connections contain a cache of HTTP header keys sent by the client. While the total number of entries in this cache is capped, an attacker sending very large keys can cause the server to allocate approximately 64 MiB per open connection.

### Affected:

| Ecosystem | Package          |
|-----------|------------------|
| Go        | stdlib           |
| Go        | golang.org/x/net |

| type   | URL                                                                                                                                                  |
|--------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| REPORT | [https://go.dev/issue/56350](https://go.dev/issue/56350)                                                                                             |
| FIX    | [https://go.dev/cl/455717](https://go.dev/cl/455717)                                                                                                 |
| FIX    | [https://go.dev/cl/455635](https://go.dev/cl/455635)                                                                                                 |
| WEB    | [https://groups.google.com/g/golang-announce/c/L_3rmdT0BMU/m/yZDrXjIiBQAJ](https://groups.google.com/g/golang-announce/c/L_3rmdT0BMU/m/yZDrXjIiBQAJ) |
## Path

Path: /scan/go.mod
Type: lockfile

### Packages

| Ecosystem | Package           | Version |
|-----------|-------------------|---------|
| Go        | golang.org/x/text | 0.3.7   |

### Details:

An attacker may cause a denial of service by crafting an Accept-Language header which ParseAcceptLanguage will take significant time to parse.

### Affected:

| Ecosystem | Package           |
|-----------|-------------------|
| Go        | golang.org/x/text |

| type   | URL                                                                                                                                                  |
|--------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| REPORT | [https://go.dev/issue/56152](https://go.dev/issue/56152)                                                                                             |
| FIX    | [https://go.dev/cl/442235](https://go.dev/cl/442235)                                                                                                 |
| WEB    | [https://groups.google.com/g/golang-announce/c/-hjNw559_tE/m/KlGTfid5CAAJ](https://groups.google.com/g/golang-announce/c/-hjNw559_tE/m/KlGTfid5CAAJ) |

# Memory and CPU log example

> Automatic function

| time                | state - running | state - dead | state - OOMKilled | state - paused | state - restarting | state - error | state - status | state - exitCode | state - health check | read      | pre read  | pids - current (linux) | pids - limit (linux) | num of process (windows) | storage - read count (windows) | storage - write count (windows) | cpu - online | cpu - system usage | cpu - usage in user mode | cpu - usage in kernel mode | cpu - total usage | cpu - throttled time | cpu - throttled periods | cpu - throttling periods | pre cpu - online | pre cpu - system usage | pre cpu - usage in user mode | pre cpu - usage in kernel mode | pre cpu - total usage | pre cpu - throttled time | pre cpu - throttled periods | pre cpu - throttling periods | memory - limit | memory - commit peak | memory - commit | memory - fail cnt | memory - usage | memory - max usage | memory - private working set |
|---------------------|-----------------|--------------|-------------------|----------------|--------------------|---------------|----------------|------------------|----------------------|-----------|-----------|------------------------|----------------------|--------------------------|--------------------------------|---------------------------------|--------------|--------------------|--------------------------|----------------------------|-------------------|----------------------|-------------------------|--------------------------|------------------|------------------------|------------------------------|--------------------------------|-----------------------|--------------------------|-----------------------------|------------------------------|----------------|----------------------|-----------------|-------------------|----------------|--------------------|------------------------------|
| 2022-12-22 18:05:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 177613586 | 166940794 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184040640000000    | 35303000                 | 39377000                   | 74681000          | 0                    | 0                       | 0                        | 8                | 184032720000000        | 32301000                     | 39031000                       | 71333000              | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 2617344        | 0                  | 0                            |
| 2022-12-22 18:05:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 171544716 | 164254632 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184118820000000    | 54666000                 | 50461000                   | 105127000         | 0                    | 0                       | 0                        | 8                | 184110970000000        | 52911000                     | 48841000                       | 101752000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 2727936        | 0                  | 0                            |
| 2022-12-22 18:05:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 171077595 | 166890387 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184196780000000    | 83039000                 | 72275000                   | 155315000         | 0                    | 0                       | 0                        | 8                | 184188900000000        | 82306000                     | 71636000                       | 153942000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 2969600        | 0                  | 0                            |
| 2022-12-22 18:05:48 | true            | false        | false             | false          | false              |               | running        |                  |                      | 263693003 | 230767753 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184284240000000    | 141004000                | 125160000                  | 266165000         | 0                    | 0                       | 0                        | 8                | 184276070000000        | 132987000                    | 118391000                      | 251378000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 3538944        | 0                  | 0                            |
| 2022-12-22 18:05:57 | true            | false        | false             | false          | false              |               | running        |                  |                      | 178493424 | 166287840 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184354520000000    | 196078000                | 181249000                  | 377327000         | 0                    | 0                       | 0                        | 8                | 184346570000000        | 189827000                    | 167112000                      | 356939000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 4067328        | 0                  | 0                            |
| 2022-12-22 18:06:07 | true            | false        | false             | false          | false              |               | running        |                  |                      | 174253762 | 166177428 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184433270000000    | 255357000                | 251859000                  | 507217000         | 0                    | 0                       | 0                        | 8                | 184425320000000        | 253246000                    | 242839000                      | 496085000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 4505600        | 0                  | 0                            |
| 2022-12-22 18:06:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 174244961 | 164435544 | 7                      | -1                   | 0                        | 0                              | 0                               | 8            | 184512140000000    | 334782000                | 312463000                  | 647245000         | 0                    | 0                       | 0                        | 8                | 184504170000000        | 325585000                    | 308840000                      | 634425000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 5455872        | 0                  | 0                            |
| 2022-12-22 18:06:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 175264632 | 164967173 | 7                      | -1                   | 0                        | 0                              | 0                               | 8            | 184590980000000    | 388605000                | 364317000                  | 752922000         | 0                    | 0                       | 0                        | 8                | 184583060000000        | 385493000                    | 359421000                      | 744915000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 5894144        | 0                  | 0                            |
| 2022-12-22 18:06:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 180010095 | 171461678 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 184669860000000    | 442193000                | 434521000                  | 876714000         | 0                    | 0                       | 0                        | 8                | 184661890000000        | 442193000                    | 427139000                      | 869332000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6246400        | 0                  | 0                            |

# Example of crash capture

> Automatic function

```
2022-12-16T04:17:02.324250972Z 2022/12/16 04:17:02 IP: 10.0.0.6
2022-12-16T04:17:02.324469972Z 2022/12/16 04:17:02 [DEBUG] memberlist: Initiating push/pull sync with:  10.0.0.6:7946
2022-12-16T04:17:02.325176930Z Member: be594d5ade2e 10.0.0.6
2022-12-16T04:17:02.325184055Z Member: 48b8e00607b2 10.0.0.7
2022-12-16T04:17:02.326104639Z 2022/12/16 04:17:02 nats connection ok
2022-12-16T04:17:02.326115514Z 2022/12/16 04:17:02 chaos enable
2022-12-16T04:17:18.253631340Z 2022/12/16 04:17:18 [DEBUG] memberlist: Stream connection from=10.0.0.8:40948
2022-12-16T04:17:21.766587550Z 2022/12/16 04:17:21 [DEBUG] memberlist: Stream connection from=10.0.0.6:54688
2022-12-16T04:17:22.327208800Z 2022/12/16 04:17:22 you can restart now
2022-12-16T04:17:26.174066386Z 2022/12/16 04:17:26 [DEBUG] memberlist: Initiating push/pull sync with: be594d5ade2e 10.0.0.6:7946
2022-12-16T04:17:33.256542333Z 2022/12/16 04:17:33 [DEBUG] memberlist: Stream connection from=10.0.0.8:48176
2022-12-16T04:17:36.770945293Z 2022/12/16 04:17:36 [DEBUG] memberlist: Stream connection from=10.0.0.6:34038
2022-12-16T04:17:41.178559545Z 2022/12/16 04:17:41 [DEBUG] memberlist: Initiating push/pull sync with: 575857e427da 10.0.0.8:7946
2022-12-16T04:17:48.258598174Z 2022/12/16 04:17:48 [DEBUG] memberlist: Stream connection from=10.0.0.8:40284
2022-12-16T04:17:56.183927386Z 2022/12/16 04:17:56 [DEBUG] memberlist: Initiating push/pull sync with: 575857e427da 10.0.0.8:7946
2022-12-16T04:18:03.260538041Z 2022/12/16 04:18:03 [DEBUG] memberlist: Stream connection from=10.0.0.8:40656
2022-12-16T04:18:03.524478167Z 2022/12/16 04:18:03 [DEBUG] memberlist: Failed UDP ping: be594d5ade2e (timeout reached)
2022-12-16T04:18:04.325072375Z 2022/12/16 04:18:04 [INFO] memberlist: Suspect be594d5ade2e has failed, no acks received
2022-12-16T04:18:05.523893459Z 2022/12/16 04:18:05 [DEBUG] memberlist: Failed UDP ping: be594d5ade2e (timeout reached)
2022-12-16T04:18:05.542607542Z 2022/12/16 04:18:05 [INFO] memberlist: Marking be594d5ade2e as failed, suspect timeout reached (1 peer confirmations)
2022-12-16T04:18:06.327539501Z 2022/12/16 04:18:06 [INFO] memberlist: Suspect be594d5ade2e has failed, no acks received
2022-12-16T04:18:11.186402295Z 2022/12/16 04:18:11 [DEBUG] memberlist: Initiating push/pull sync with: 575857e427da 10.0.0.8:7946
2022-12-16T04:18:18.262939049Z 2022/12/16 04:18:18 [DEBUG] memberlist: Stream connection from=10.0.0.8:34016
2022-12-16T04:19:38.808491836Z 2022/12/16 04:19:38 IP: 10.0.0.6
2022-12-16T04:19:38.808835169Z 2022/12/16 04:19:38 [DEBUG] memberlist: Initiating push/pull sync with:  10.0.0.6:7946
2022-12-16T04:19:38.809361419Z Member: be594d5ade2e 10.0.0.6
2022-12-16T04:19:38.809370169Z Member: 48b8e00607b2 10.0.0.7
2022-12-16T04:19:40.810843628Z 2022/12/16 04:19:40 nats connection error: read tcp 10.0.0.7:39924->10.0.0.2:4222: i/o timeout
2022-12-16T04:19:40.810940878Z 2022/12/16 04:19:40 bug: messageSystem.Subscribe().error: nats: invalid connection
```