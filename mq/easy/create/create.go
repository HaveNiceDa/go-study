package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	flag.Parse()
	queues := flag.Args()
	if len(queues) == 0 {
		log.Fatal("必须指定队列名称，例如: go run create.go consumer_a consumer_b")
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

	for _, name := range queues {
		q, err := ch.QueueDeclare(
			name,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("声明队列 %s 失败: %s", name, err)
		}

		err = ch.QueueBind(
			q.Name,
			"",
			exChangeName,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("绑定队列 %s 失败: %s", name, err)
		}
		log.Printf("队列 %s 已声明并绑定到 %s", name, exChangeName)
	}

	for i := 0; i < 20; i++ {
		body := fmt.Sprintf("msg-%d", i)
		err = ch.Publish(
			exChangeName,
			"",
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
		fmt.Printf("发送消息: %s\n", body)
	}
}
