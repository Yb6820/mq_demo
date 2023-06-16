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
	//定义交换机的名称
	exchangeName := "Direct_Exchange"
	//定义队列的名称
	queueNames := []string{"direct_Queue1", "direct_Queue2", "direct_Queue3", "direct_Queue4"}
	//定义Key值
	keys := []string{"key_1", "key_3", "key_4"}

	//申请通道
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	//声明队列
	q1, err := ch.QueueDeclare(queueNames[0], true, false, false, false, nil)
	FailOnError(err, "Failed to Create a Queue")
	q2, err := ch.QueueDeclare(queueNames[1], true, false, false, false, nil)
	FailOnError(err, "Failed to Create a Queue")
	q3, err := ch.QueueDeclare(queueNames[2], true, false, false, false, nil)
	FailOnError(err, "Failed to Create a Queue")
	q4, err := ch.QueueDeclare(queueNames[3], true, false, false, false, nil)
	FailOnError(err, "Failed to Create a Queue")

	//声明交换机类型
	err = ch.ExchangeDeclare(
		exchangeName, //交换机的名称
		//交换机的类型，分为：direct(直连),fanout(扇出,类似广播),topic(话题,与direct相似但是模式匹配),headers(用header来设置生产和消费的key)
		"direct",
		true,  //是否持久化
		false, //是否自动删除
		false, //是否公开，false即公开
		false, //是否等待
		nil,
	)
	FailOnError(err, "Failed to Declare a Exchange")

	//根据key将队列与keys绑定
	ch.QueueBind(q1.Name, keys[0], exchangeName, false, nil)
	ch.QueueBind(q2.Name, keys[0], exchangeName, false, nil)
	ch.QueueBind(q3.Name, keys[1], exchangeName, false, nil)
	ch.QueueBind(q4.Name, keys[2], exchangeName, false, nil)

	//发送消息
	err = ch.Publish(exchangeName, keys[0], false, false,
		amqp.Publishing{
			Type: "text/plain",
			Body: []byte("Hello Dierct key1 message"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeName, "Hello Dierct key1 message")
	FailOnError(err, "Failed to publish a message")

	err = ch.Publish(exchangeName, keys[1], false, false,
		amqp.Publishing{
			Type: "text/plain",
			Body: []byte("Hello Dierct key3 message"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeName, "Hello Dierct key3 message")
	FailOnError(err, "Failed to publish a message")

	err = ch.Publish(exchangeName, keys[2], false, false,
		amqp.Publishing{
			Type: "text/plain",
			Body: []byte("Hello Dierct key4 message"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeName, "Hello Dierct key4 message")
	FailOnError(err, "Failed to publish a message")
}
