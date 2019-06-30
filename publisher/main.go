package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"

	"github.com/phamvinhdat/rabbitmq/share"
	"github.com/phamvinhdat/rabbitmq/utils"
)

func main() {
	conn, err := utils.NewConnectRabbitMQ(share.NewConfiguration().AMQPConnectionURL)
	utils.HandleError(err, "Can not connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	utils.HandleError(err, "Can not create a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	utils.HandleError(err, "Can not declare 'add' queue")

	for {
		rand.Seed(time.Now().UnixNano())
		addTask := share.AddTask{
			Number1: rand.Intn(1000),
			Number2: rand.Intn(1000),
		}
		body, err := json.Marshal(addTask)
		if err != nil {
			utils.HandleError(err, "Error encoding JSON")
		}

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

		if err != nil {
			log.Fatalln("Error publishing message:", err)
		}

		log.Printf("AddTask: %d + %d", addTask.Number1, addTask.Number2)

		time.Sleep(time.Millisecond * 750)
	}
}
