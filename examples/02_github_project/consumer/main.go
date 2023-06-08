package main

import (
	"bytes"
	"fmt"
	polar "github.com/polarstreams/go-client"
	"io/fs"
	"os"
)

type Order struct {
	ID    string
	Item  string
	Price float64
	EMail string
}

func main() {
	var err error
	var consumer polar.Consumer
	consumer, err = NewConsumer("polar://delete_polar_0#polar://delete_polar_1#polar://delete_polar_2")
	if err != nil {
		panic(err)
	}

	fmt.Println("receiving...")
	err = ReceiveData(consumer)
	if err != nil {
		panic(err)
	}
}

func SaveData(data []byte) (err error) {
	var f *os.File
	f, err = os.OpenFile("/data/ignore.dataReceived.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, fs.ModePerm)
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

func ReceiveData(consumer polar.Consumer) (err error) {
	for {
		pollResult := consumer.Poll()
		if pollResult.Error != nil {
			fmt.Printf("Found error while polling: %s", pollResult.Error)
			continue
		}

		// New records organized by topic
		for _, topicRecords := range pollResult.TopicRecords {
			for _, record := range topicRecords.Records {
				fmt.Println(string(record.Body), record.Timestamp)

				if bytes.Equal(record.Body, []byte("{\"end\":\"end\"}")) {
					return
				}

				err = SaveData(record.Body)
				if err != nil {
					return
				}
			}
		}
	}
}

func NewConsumer(url string) (consumer polar.Consumer, err error) {
	topic := "topic" // The topic will be automatically created
	group := "group1"
	consumer, err = polar.NewConsumer(url, group, topic)
	if err != nil {
		return
	}

	fmt.Printf("Discovered a cluster with %d brokers\n", consumer.BrokersLength())

	return
}
