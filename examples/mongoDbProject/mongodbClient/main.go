package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"runtime/debug"
	"time"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Add(-3*time.Hour).Format("2006/01/02 15:04:05") + " " + string(bytes))
}

func main() {
	var err error
	var mongoClient *mongo.Client
	var start = time.Now()

	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	log.Println("container is running")

	// "mongodb://delete_delay_0:27017,delete_delay_1:27017,delete_delay_2:27017/?replicaSet=rs0"
	connectingString := os.Getenv("CONNECTION_STRING")

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(connectingString))
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Connect(context.Background())
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Printf("error: %v\n", err.Error())
		panic(string(debug.Stack()))
	}

	type Trainer struct {
		Name string
		Age  int
		City string
	}

	var collection *mongo.Collection
	var totalOfInserts int64 = 10000000
	for i := int64(0); i != totalOfInserts; i += 1 {
		collection = mongoClient.Database("test").Collection("trainers")
		trainer := Trainer{gofakeit.Name(), gofakeit.Number(18, 99), gofakeit.City()}
		_, err = collection.InsertOne(context.Background(), trainer)
		if err != nil {
			panic(err)
		}

		if i%1000 == 0 {
			log.Printf("inserted %v", i)
		}
	}

	var total int64
	if total, err = collection.CountDocuments(context.Background(), bson.M{}); err != nil {
		panic(err)
	}

	if total != totalOfInserts {
		err = fmt.Errorf("total of inserts must be %v found %v", totalOfInserts, total)
		panic(err)
	}

	log.Printf("end of test\n")
	duration := time.Since(start)
	log.Printf("Duration: %v\n\n", duration)
}
