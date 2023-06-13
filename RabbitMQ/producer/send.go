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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//创建RabbitMQ中的管道
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//声明队列，如果队列不存在则自动创建，存在则跳过创建
	q, err := ch.QueueDeclare(
		"hello", //消息队列的名称
		false,   //是否持久化
		false,   //是否自动删除
		false,   //是否具有排他性
		false,   //是否阻塞处理
		nil,     //额外的属性
	)
	FailOnError(err, "Failed to declare a queue")

	body := "Hello World"
	err = ch.Publish(
		"",
		q.Name,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	log.Printf(" [x] Sent %s", body)
	FailOnError(err, "Failed to publish a message")
}
