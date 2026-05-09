package main

import (
	"fmt"
	"log"
	"os"
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

	exChangeName := "logs_topic"

	err = ch.ExchangeDeclare(
		exChangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换器失败 %s", err)
	}

	bindings := []string{"kern.*", "auth.*"}
	for _, b := range bindings {
		q, err := ch.QueueDeclare(
			"queue_"+strings.ReplaceAll(b, "*", "all"),
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
			b,
			exChangeName,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("绑定队列失败 %s", err)
		}
		log.Printf("队列 %s 已绑定到 %s (binding key: %s)", q.Name, exChangeName, b)
	}

	messages := []struct {
		routingKey string
		body       string
	}{
		{routingKey: "kern.critical", body: "内核严重错误"},
		{"kern.info", "内核信息"},
		{"auth.warning", "认证警告"},
		{"auth.info", "认证信息"},
		{"cron.critical", "定时任务严重错误"},
	}

	for _, m := range messages {
		err = ch.Publish(
			exChangeName,
			m.routingKey,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(m.body),
			})
		if err != nil {
			log.Printf("发送失败: %v", err)
			return
		}
		fmt.Printf("[%s] 发送消息: %s\n", m.routingKey, m.body)
	}

	if len(os.Args) > 2 {
		routingKey := os.Args[1]
		body := strings.Join(os.Args[2:], " ")
		err = ch.Publish(
			exChangeName,
			routingKey,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		if err != nil {
			log.Printf("发送失败: %v", err)
			return
		}
		fmt.Printf("[%s] 发送消息: %s\n", routingKey, body)
	}
}
