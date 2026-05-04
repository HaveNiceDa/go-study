package main

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	// 连接 RabbitMQ
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法打开通道: %v", err)
	}
	defer ch.Close()

	queueName := "xxx"

	// 声明队列
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		false,     // 持久性
		false,     // 自动删除
		false,     // 排他性
		false,     // 非阻塞
		nil,       // 其他参数
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	// 发送消息
	body := fmt.Sprintf("msg-%d", 1)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Body: []byte(body),
		})
	if err != nil {
		logrus.Errorf("发送失败: %v", err)
		return
	}
	fmt.Printf("发送消息: %s\n", body)
}
