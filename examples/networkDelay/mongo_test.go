package mongodbproject

import (
	"context"
	"fmt"
	"github.com/helmutkemper/chaos/factory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime/debug"
	"testing"
	"time"
)

func TestLinearNetworkWithProblems(t *testing.T) {

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		Test(t)

	factory.NewContainerFromImage(
		"mongo:latest",
	).
		// Limita a origem de conexão ao MongoDB
		EnvironmentVar([]string{"bindIp:delete_delay_0"}).
		Create("mongo", 1).
		Start()

	// Cria um container com a propriedade de interromper pacotes de rede e simular uma rede com problemas
	factory.NewContainerNetworkProxy(
		"delay",

		// Porta local é a porta que sofre conexão externa
		27016,

		// Destino é sempre o elemento passivo. Por exemplo, o banco de dados sempre recebe conexões, é passivo.
		// Cuidado quando usar host name, ele só funciona na rede docker quando é criada uma nova rede
		"delete_mongo_0:27017",

		// O atraso é em milissegundo. Lembre-se que a rede transporta pacotes e cada pacote recebe um valor aleatório entre o mínimo e o máximo.
		10, 100,
	)

	// Lembre-se, a porta 27017 é a porta original do banco e tem acesso normal, a porta 27016 é a porta da rede com problemas
	go func() {
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
	}()

	if !primordial.Monitor(2 * time.Minute) {
		t.Fail()
	}
}
