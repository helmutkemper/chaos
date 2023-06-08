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
		//Ports("tcp", 9250, 9250).                  |
		//Ports("tcp", 9251, 9251).          +-------+
		//Ports("tcp", 9252, 9252).          |
		//                                   ↓
		// Determines the name of the container and the number of containers to be created
		Create("polar", 3).
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
