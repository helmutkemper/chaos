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
