package docker

import (
	"context"
	"fmt"
	dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
	"github.com/helmutkemper/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func ExampleNetworkChaos_Init() {
	var err error

	SaGarbageCollector()

	var netDocker = &dockerNetwork.ContainerBuilderNetwork{}
	err = netDocker.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	var mongoDocker = &ContainerBuilder{}
	mongoDocker.SetNetworkDocker(netDocker)
	mongoDocker.SetImageName("mongo:latest")
	mongoDocker.SetContainerName("container_delete_mongo_after_test")
	//mongoDocker.AddPortToChange("27017", "27016")
	//mongoDocker.AddPortToExpose("27017")
	mongoDocker.SetEnvironmentVar(
		[]string{
			"--bind_ip_all",
			"--host 0.0.0.0",
			"--bind 0.0.0.0",
		},
	)
	mongoDocker.SetWaitStringWithTimeout(`"msg":"Waiting for connections","attr":{"port":27017`, 20*time.Second)
	err = mongoDocker.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	err = mongoDocker.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	var redis = ContainerBuilder{}
	redis.SetNetworkDocker(netDocker)
	redis.SetImageName("redis:latest")
	redis.SetContainerName("container_delete_redis_test")
	redis.AddPortToExpose("6379")
	redis.SetWaitStringWithTimeout("Ready to accept connections", 10*time.Second)

	err = redis.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	err = redis.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	var chaos = NetworkChaos{}
	chaos.SetNetworkDocker(netDocker)
	chaos.SetFatherContainer(mongoDocker)
	chaos.SetPorts(27017, 27016, false)
	err = chaos.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	err = testNetworkOverloaded(
		"mongodb://0.0.0.0:27016",
		2*time.Second,
	)

	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	//Output:
	//
}

// testNetworkOverloaded (English): Tests the new network port
// testNetworkOverloaded (Português): Testa a nova porta de rede
func testNetworkOverloaded(
	address string,
	timeout time.Duration,
) (
	err error,
) {

	// (English): Runtime measurement starts
	// (Português): Começa a medição do tempo de execução
	start := time.Now()

	var mongoClient *mongo.Client
	var cancel context.CancelFunc
	var ctx context.Context

	// (English): Prepare the MongoDB client
	// (Português): Prepara o cliente do MongoDB
	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		return
	}

	// (English): Connects to MongoDB
	// (Português): Conecta ao MongoDB
	err = mongoClient.Connect(ctx)
	if err != nil {
		return
	}

	// (English): Prepares the timeout
	// (Português): Prepara o tempo limite
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// (English): Ping() to test the MongoDB connection
	// (Português): Faz um ping() para testar a conexão do MongoDB
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	// (English): New collection format
	// (Português): Formato da nova coleção
	type Trainer struct {
		Name string
		Age  int
		City string
	}

	// (English): Creates the 'test' bank and the 'dinos' collection
	// (Português): Cria o banco 'test' e a coleção 'dinos'
	collection := mongoClient.Database("test").Collection("dinos")

	// (English): Prepares the data to be inserted
	// (Português): Prepara os dados a serem inseridos
	trainerData := Trainer{"T-Rex", 10, "Jurassic Town"}

	for i := 0; i != 100; i += 1 {
		// (English): Insert the data
		// (Português): Insere os dados
		_, err = collection.InsertOne(context.TODO(), trainerData)
		if err != nil {
			log.Printf("collection.InsertOne().error: %v", err)
			return
		}
	}

	// (English): Stop the operation time measurement
	// (Português): Para a medição de tempo da operação
	duration := time.Since(start)
	fmt.Printf("End!\n")
	fmt.Printf("Duration: %v\n\n", duration)

	return
}
