package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	routingKey := flag.String("key", "", "路由键（node001 或 node002）")
	flag.Parse()
	if *routingKey == "" {
		log.Fatal("必须指定 -key 参数，例如: go run consumer.go -key node001")
	}

	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法打开通道: %v", err)
	}
	defer ch.Close()

	exChangeName := "logs_direct"

	err = ch.ExchangeDeclare(
		exChangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换器失败 %s", err)
	}

	q, err := ch.QueueDeclare(
		"queue_"+*routingKey,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败 %s", err)
	}

	err = ch.QueueBind(
		q.Name,
		*routingKey,
		exChangeName,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定队列失败 %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("创建消费者失败 %s", err)
	}

	fmt.Printf("[%s] 准备接收消息\n", *routingKey)
	for d := range msgs {
		log.Printf("[%s] 收到消息: %s\n", *routingKey, d.Body)
	}
}
