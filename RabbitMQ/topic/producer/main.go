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
	//建立连接
	conn, err := amqp.Dial("amqp://youbet:181234@localhost/guest")
	FailOnError(err, "Failed to Connected to RabbitMQ")
	defer conn.Close()

	//定义交换机的名称
	exchangeName := "Topic_Exchange"
	//定义队列的名称
	queueNames := []string{"topic_Queue1", "topic_Queue2", "topic_Queue3", "topic_Queue4"}
	//定义keys
	keys := []string{"key1.key2.key3.*", "key1.#", "*.key2.*.key4", "#.key3.key4"}

	//申请通道
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//声明队列
	q1, err := ch.QueueDeclare(queueNames[0], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")
	q2, err := ch.QueueDeclare(queueNames[1], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")
	q3, err := ch.QueueDeclare(queueNames[2], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")
	q4, err := ch.QueueDeclare(queueNames[3], true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")

	//声明交换机
	err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
	FailOnError(err, "Failed to Declare a Exchange")

	//将队列和key绑定到交换机上
	ch.QueueBind(q1.Name, keys[0], exchangeName, false, nil)
	ch.QueueBind(q2.Name, keys[1], exchangeName, false, nil)
	ch.QueueBind(q3.Name, keys[2], exchangeName, false, nil)
	ch.QueueBind(q4.Name, keys[3], exchangeName, false, nil)

	//发送消息
	err = ch.Publish(exchangeName, "key1.key2.key3.key4", false, false,
		amqp.Publishing{
			Type: "text/plain",
			Body: []byte("Hello Topic key1 message"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeName, "Hello Topic key1 message")
	FailOnError(err, "Failed to publish a message")
}
