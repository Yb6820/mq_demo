package main

import (
	"log"

	"github.com/streadway/amqp"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s:%s", msg, err)
	}
}
func main() {
	conn, err := amqp.Dial("amqp://youbet:181234@localhost:5672/guest")
	FailOnError(err, "Failed to Connected to RabbitMQ")
	defer conn.Close()
	//定义队列的名称
	queueName := "Headers_Queue1"

	//申请通道
	ch, err := conn.Channel()
	FailOnError(err, "Failed to Open a Channel")
	defer ch.Close()

	//消费消息
	forever := make(chan bool)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)

	FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message From %s : %s", queueName, d.Body)
		}
	}()
	<-forever
}
