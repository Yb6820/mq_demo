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
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//创建RabbitMQ中的管道
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//为消息队列注册消费者
	//接收消息
	forever := make(chan bool)
	msgs, err := ch.Consume(
		"hello", // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	FailOnError(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
