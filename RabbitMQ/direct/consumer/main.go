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
	exchangeNames := []string{"direct_exchange1", "direct_exchange3", "direct_exchange4"}

	queueNames := []string{"direct_Queue1", "direct_Queue2", "direct_Queue3", "direct_Queue4"}

	//获取一个通道
	ch, err := conn.Channel()
	FailOnError(err, "Failed to Create a channel")
	defer ch.Close()

	q1, err := ch.QueueDeclare(queueNames[0], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")

	q2, err := ch.QueueDeclare(queueNames[1], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")

	q3, err := ch.QueueDeclare(queueNames[2], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")

	q4, err := ch.QueueDeclare(queueNames[3], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")

	//消费者绑定队列
	//direct模式下的key与要订阅的虚拟机的key保持一直
	ch.QueueBind(q1.Name, q1.Name, exchangeNames[0], false, nil)

	ch.QueueBind(q2.Name, q1.Name, exchangeNames[0], false, nil)

	ch.QueueBind(q3.Name, q3.Name, exchangeNames[1], false, nil)

	ch.QueueBind(q4.Name, q4.Name, exchangeNames[2], false, nil)

	//消费消息
	forever := make(chan bool)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	msgs, err := ch.Consume(q1.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message From %s : %s", queueNames[0], d.Body)
		}
	}()

	msgs, err = ch.Consume(q2.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message From %s : %s", queueNames[1], d.Body)
		}
	}()

	msgs, err = ch.Consume(q3.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message From %s : %s", queueNames[2], d.Body)
		}
	}()

	msgs, err = ch.Consume(q4.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message From %s : %s", queueNames[3], d.Body)
		}
	}()
	<-forever
}
