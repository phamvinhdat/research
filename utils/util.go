package utils

import (
	"github.com/streadway/amqp"
	"log"
)

func HandleError(err error, msg string){
	if err != nil{
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewConnectRabbitMQ(amqpConnectionURL string)(*amqp.Connection, error){
	conn, err := amqp.Dial(amqpConnectionURL)
	if err != nil{
		return nil, err
	}
	return conn, nil
}