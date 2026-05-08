package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	queueName := flag.String("queue", "", "队列名称（每个消费者使用不同名称）")
	flag.Parse()
	if *queueName == "" {
		log.Fatal("必须指定 -queue 参数，例如: go run consumer.go -queue consumer_a")
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

	exChangeName := "logs_persistent"

	err = ch.ExchangeDeclare(
		exChangeName,
		"fanout",
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
		*queueName,
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
		"",
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

	fmt.Printf("[%s] 准备接收消息\n", *queueName)
	for d := range msgs {
		log.Printf("[%s] 收到消息: %s\n", *queueName, d.Body)
	}
}
