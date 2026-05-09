package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
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

	msgs, err := ch.Consume(
		"dead_letter_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("创建死信消费者失败 %s", err)
	}

	fmt.Println("死信消费者启动，等待死信消息...")

	for d := range msgs {
		log.Printf("[死信] 收到死信消息: %s | 原始routing-key: %s\n", d.Body, d.RoutingKey)
		d.Ack(false)
	}
}
