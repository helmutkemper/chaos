package manager

import (
	"context"
	"github.com/helmutkemper/chaos/internal/standalone"
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestContainerFromImage_Primordial(t *testing.T) {
	var errorCh = make(chan error)
	go func(t *testing.T) {
		err := <-errorCh
		t.Error(err.Error())
		t.Fail()
	}(t)

	standalone.GarbageCollector()
	t.Cleanup(func() {
		//standalone.GarbageCollector()
	})

	mng := &Manager{}
	mng.New(errorCh)

	mng.Primordial().
		NetworkCreate("delete_before_test", "10.0.0.0/16", "10.0.0.1")
	mng.ContainerFromImage().
		SaveStatistics("../../").
		Ports("tcp", 27017, 27016, 27015, 27014).
		//Volumes("/data/db", "../../internal/builder/test/data0", "../../internal/builder/test/data1", "../../internal/builder/test/data2").
		EnvironmentVar("--host 0.0.0.0").
		Create("mongo:latest", "delete_mongo", 3).
		Start()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27016"))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	usersCollection := client.Database("testing").Collection("users")

	user := bson.D{{"fullName", "User 1"}, {"age", 30}}
	// insert the bson object using InsertOne()
	for {
		_, err = usersCollection.InsertOne(context.TODO(), user)
		// check for errors in the insertion
		if err != nil {
			panic(err)
		}
	}

	done := make(chan struct{})
	done <- struct{}{}
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
