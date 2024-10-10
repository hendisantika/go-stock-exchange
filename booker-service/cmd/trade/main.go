package main

import (
	"booker-service/internal/infra/kafka"
	"booker-service/internal/market/dto"
	"booker-service/internal/market/entity"
	"booker-service/internal/market/transformer"
	"encoding/json"
	"fmt"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)

	kafkaMsgChannel := make(chan *ckafka.Message)

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest", //only get messages that were created after the program was uploaded
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	//create a new thread to run the consumer without blocking the execution of the rest of the program
	go kafka.Consume(kafkaMsgChannel)

	//receives from kafka, plays on the input channel, processes, plays on the output channel and publishes on kafka
	book := entity.NewBook(ordersIn, ordersOut, wg)

	//create new thread to process orders
	go book.Trade()

	//anonymous function to receive messages from kafka
	go func() {
		for msg := range kafkaMsgChannel {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			//deserialize the json message -> dto
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)

			if err != nil {
				panic(err)
			}

			order := transformer.TransformInput(tradeInput)

			//throws orders to the book input channel
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		orderOutput := transformer.TransformOutput(res)
		outputJson, err := json.Marshal(orderOutput)

		fmt.Print("Negotiation completed:")
		fmt.Println(string(outputJson))

		if err != nil {
			fmt.Println(err)
		}

		//publish the processed orders back to kafka
		producer.Publish(outputJson, []byte("orders"), "output")

	}
}
