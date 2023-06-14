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
	exchangeNames := []string{"direct_exchange1", "direct_exchange3", "direct_exchange4"}
	//定义队列的名称
	queueNames := []string{"direct_Queue1", "direct_Queue2", "direct_Queue3", "direct_Queue4"}

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

	//声明交换机
	err = ch.ExchangeDeclare(
		exchangeNames[0], //交换机的名称
		//交换机的类型，分为：direct(直连),fanout(扇出,类似广播),topic(话题,与direct相似但是模式匹配),headers(用header来设置生产和消费的key)
		"direct",
		true,  //是否持久化
		false, //是否自动删除
		false, //是否公开，false即公开
		false, //是否等待
		nil,
	)
	FailOnError(err, "Failed to Declare a Exchange")

	err = ch.ExchangeDeclare(exchangeNames[1], "direct", true, false, false, false, nil)
	FailOnError(err, "Failed to Declare a Exchange")

	err = ch.ExchangeDeclare(exchangeNames[2], "direct", true, false, false, false, nil)
	FailOnError(err, "Failed to Declare a Exchange")

	//根据key将队列与交换机绑定
	ch.QueueBind(q1.Name, q1.Name, exchangeNames[0], false, nil)
	ch.QueueBind(q2.Name, q1.Name, exchangeNames[0], false, nil)
	ch.QueueBind(q3.Name, q3.Name, exchangeNames[1], false, nil)
	ch.QueueBind(q4.Name, q4.Name, exchangeNames[2], false, nil)

	//发送消息
	body := "Hello DirectExchange"
	err = ch.Publish(exchangeNames[0], q1.Name, false, false,
		amqp.Publishing{
			Type: "text/plain",
			Body: []byte(body + "1"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeNames[0], body+"1")
	FailOnError(err, "Failed to publish a message")

	err = ch.Publish(exchangeNames[1], q3.Name, false, false,
		amqp.Publishing{
			Type: "text/plain",
			Body: []byte(body + "3"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeNames[1], body+"3")
	FailOnError(err, "Failed to publish a message")

	err = ch.Publish(exchangeNames[2], q4.Name, false, false,
		amqp.Publishing{
			Type: "text/plain",
			Body: []byte(body + "4"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeNames[2], body+"4")
	FailOnError(err, "Failed to publish a message")
}
