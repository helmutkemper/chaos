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
// This example will show the basic settings of how to create an image-based container and how to expose a door to the
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
