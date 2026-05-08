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

	routingKeys := []string{"node001", "node002"}
	for _, key := range routingKeys {
		q, err := ch.QueueDeclare(
			"queue_"+key,
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
			key,
			exChangeName,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("绑定队列失败 %s", err)
		}
		log.Printf("队列 queue_%s 已绑定到 %s (routing key: %s)", key, exChangeName, key)
	}

	for i := 0; i < 20; i++ {
		body := fmt.Sprintf("msg-%d", i)
		err = ch.Publish(
			exChangeName,
			"node001",
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         []byte(body),
			})
		if err != nil {
			log.Printf("发送失败: %v", err)
			return
		}
		fmt.Printf("[node001] 发送消息: %s\n", body)
	}
	for i := 0; i < 10; i++ {
		body := fmt.Sprintf("msg-%d", i)
		err = ch.Publish(
			exChangeName,
			"node002",
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         []byte(body),
			})
		if err != nil {
			log.Printf("发送失败: %v", err)
			return
		}
		fmt.Printf("[node002] 发送消息: %s\n", body)
	}
}
