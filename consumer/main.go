package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

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

	err = amqpChannel.Qos(1, 0, false)
	utils.HandleError(err, "Can not configure Qos")

	messageChannel, err := amqpChannel.Consume(queue.Name, "", false, false, false, false, nil)
	utils.HandleError(err, "Can not register consumer")

	stopChan := make(chan bool)
	go func() {
		log.Println("Consumer ready. PID: ", os.Getegid())
		for d := range messageChannel {
			log.Println("Received a message:", string(d.Body))

			//processing
			time.Sleep(time.Second * 1)

			addTask := &share.AddTask{}
			err := json.Unmarshal(d.Body, addTask)
			if err != nil {
				log.Println("Error decoding JSON:", err)
			}

			_handlerTask(addTask, func() {
				log.Printf("%d + %d = %d\n", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)
			})

			if err = d.Ack(false); err != nil {
				log.Println("Error acknowledging message:", err)
			} else {
				log.Println("Acknowledged message")
			}
		}
	}()

	//stop for program termination
	<-stopChan
}

func _handlerTask(input interface{}, f func()) {
	f()
}
