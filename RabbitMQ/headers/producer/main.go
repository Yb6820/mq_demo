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
	exchangeName := "Headers_Exchange"
	//定义队列的名称
	queueName := "Headers_Queue1"

	//申请通道
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//声明队列
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	FailOnError(err, "Failed to declare a Queue")

	//声明交换机
	err = ch.ExchangeDeclare(exchangeName, "headers", true, false, false, false, nil)
	FailOnError(err, "Failed to Declare a Exchange")

	//绑定交换机与队列的headers对应关系
	headers := make(map[string]interface{})
	headers["x-match"] = "all"
	headers["name"] = "zhangsan"
	headers["sex"] = "男"
	//最后一个参数绑定headers
	ch.QueueBind(queueName, "", exchangeName, false, headers)

	//发送消息
	err = ch.Publish(exchangeName, "", false, false,
		amqp.Publishing{
			//设置发送消息时的headers信息
			Headers: amqp.Table{
				"name": "zhangsan",
				"sex":  "男",
			}, //这里把刚才定义的headers添加进来
			Type: "text/plain",
			Body: []byte("Hello ALL Headers message"),
		},
	)
	log.Printf(" [x] Sent to %s : %s", exchangeName, "Hello ALL Headers message")
	FailOnError(err, "Failed to publish a message")
}
