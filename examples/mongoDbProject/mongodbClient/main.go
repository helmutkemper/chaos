package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime/debug"
	"time"
)

func main() {
	var err error
	var mongoClient *mongo.Client
	var start = time.Now()

	fmt.Println("container is running")

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://delete_mongo_0:27017,delete_mongo_1:27017,delete_mongo_2:27017/?replicaSet=rs0"))
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
		err = fmt.Errorf("total of inserts must be %v found %v", totalOfInserts, total)
		panic(err)
	}

	fmt.Printf("end of test\n")
	duration := time.Since(start)
	fmt.Printf("Duration: %v\n\n", duration)
}
