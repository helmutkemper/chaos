package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/polarstreams/go-client"
	"io/fs"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
)

type Order struct {
	ID    string
	Item  string
	Price float64
	EMail string
}

func (el *Order) NewFake() {
	el.ID = gofakeit.UUID()
	el.Item = gofakeit.BeerAlcohol()
	el.Price = gofakeit.Price(0.99, 10.99)
	el.EMail = gofakeit.Email()
}

func main() {
	var err error
	err = Event()
	if err != nil {
		panic(err)
	}
}

func Event() (err error) {
	var producer polar.Producer
	producer, err = NewProducer("polar://delete_polar_0")
	if err != nil {
		return
	}

	for {
		var order Order
		order.NewFake()

		var data []byte
		data, err = json.Marshal(&order)
		if err != nil {
			return
		}

		err = SaveData(data)
		if err != nil {
			return
		}

		err = SendData(data, producer)
		if err != nil {
			return
		}

		fmt.Println(string(data))

		time.Sleep(1 * time.Second)

		_, err = os.ReadFile("/data/ignore.end.empty")
		if err == nil {
			_ = os.Remove("/data/ignore.end.empty")
			_ = SendData([]byte("{\"end\":\"end\"}"), producer)
			return
		}
	}
}

func SaveData(data []byte) (err error) {
	var f *os.File
	f, err = os.OpenFile("/data/ignore.dataSent.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, fs.ModePerm)
	if err != nil {
		return
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return
	}

	_, err = f.WriteString("\n")
	return
}

func SendData(data []byte, producer polar.Producer) (err error) {

	topic := "topic" // The topic will be automatically created
	message := bytes.NewReader(data)
	partitionKey := "" // Empty to use a random partition

	return producer.Send(topic, message, partitionKey)
}

// "polar://delete_polar_0"

func NewProducer(url string) (producer polar.Producer, err error) {
	producer, err = polar.NewProducer(url)
	if err != nil {
		return
	}

	fmt.Printf("Discovered a PolarStreams cluster with %d brokers\n", producer.BrokersLength())

	return
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
//
//
//
//
//
//
//
