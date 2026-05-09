package main

import (
	"fmt"
	"log"
	"strings"

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
		"normal_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("创建消费者失败 %s", err)
	}

	fmt.Println("正常消费者启动，等待消息...")

	for d := range msgs {
		body := string(d.Body)
		if strings.Contains(body, "失败") {
			log.Printf("[拒绝] 消息: %s → 转入死信队列\n", body)
			d.Nack(false, false)
		} else {
			log.Printf("[消费] 消息: %s\n", body)
			d.Ack(false)
		}
	}
}
