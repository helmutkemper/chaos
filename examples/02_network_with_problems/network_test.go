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

// Este é um teste com simulação de rede com problemas
// Caso tenha pulado a explicação anterior, ela contém o conhecimento básico de uso do sistema. Aqui são adicionadas mais informações
//
// Neste exemplo será mostrado como criar um container com a capacidade de atrapalhar a conexão de rede
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
		"mongo:latest",
	).
		// Limita a origem de conexão ao MongoDB
		// Como o container de simulação de rede tem o nome "delay", o container será criado com nome, e o host name, "delete_delay_0"
		EnvironmentVar([]string{"bindIp:delete_delay_0"}).
		Create("mongo", 1).
		Start()

	// Cria um container com a propriedade de interromper pacotes de rede e simular uma rede com problemas
	factory.NewContainerNetworkProxy(
		"delay",

		// Uma configuração para cada container proxy
		[]factory.ProxyConfig{
			{
				// Porta de entrada do mundo externo
				LocalPort: 27016,
				// Conexão com elemento passivo, nesse caso, o MongoDB
				Destination: "delete_mongo_0:27017",

				// Tempo mínimo e máximo para atraso entre pacotes
				// total test time: ~1.8s
				//MinDelay: 0,
				//MaxDelay: 0,

				// Tempo mínimo e máximo para atraso entre pacotes
				// total test time: ~1m58
				MinDelay: 100,
				MaxDelay: 130,

				// Tempo mínimo e máximo para atraso entre pacotes
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
