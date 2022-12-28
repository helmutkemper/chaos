package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	networkdelay "github.com/helmutkemper/chaos/internal/networkDelay"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime/debug"
	"sync"
	"time"
)

var localAddr *string = flag.String("l", ":27016", "local address")
var remoteAddr *string = flag.String("r", ":27017", "remote address")

type ParserFunc struct{}

func (e ParserFunc) Parser(data []byte, direction string) (dataSize int, err error) {
	fmt.Printf("direction: %v\n%v\n", direction, hex.Dump(data))
	dataSize = len(data)
	return
}

func main() {
	flag.Parse()
	fmt.Printf("Listening: %v\nProxying: %v\n\n", *localAddr, *remoteAddr)

	var p networkdelay.ParserInterface = &ParserFunc{}

	var proxy networkdelay.Proxy
	proxy.SetBufferSize(32 * 1024)
	proxy.SetParserFunction(p)
	proxy.SetDelayMillesecond(500, 1000)
	go func() {
		var err error
		err = proxy.Proxy(*localAddr, *remoteAddr)
		if err != nil {
			panic(err)
		}
	}()

	var err error
	var timeOut = time.Second * 15

	var mongoClient *mongo.Client
	var ctx context.Context

	start := time.Now()

	fmt.Printf("conexão\n")

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27016"))
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Connect(ctx)
	if err != nil {
		panic(string(debug.Stack()))
	}

	ctx, _ = context.WithTimeout(context.Background(), timeOut)
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		panic(string(debug.Stack()))
	}

	type Trainer struct {
		Name string
		Age  int
		City string
	}

	for i := 0; i != 100; i += 1 {
		collection := mongoClient.Database("test").Collection("trainers")
		ash := Trainer{"Ash", 10, "Pallet Town"}
		_, err = collection.InsertOne(context.TODO(), ash)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("fim\n")
	duration := time.Since(start)
	fmt.Printf("Duration: %v\n\n", duration)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
