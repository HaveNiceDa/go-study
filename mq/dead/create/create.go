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

	// 声明死信交换器
	err = ch.ExchangeDeclare(
		"dlx_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明死信交换器失败 %s", err)
	}

	// 声明死信队列
	_, err = ch.QueueDeclare(
		"dead_letter_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明死信队列失败 %s", err)
	}

	// 绑定死信队列到死信交换器
	err = ch.QueueBind(
		"dead_letter_queue",
		"dead_letter_routing_key",
		"dlx_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定死信队列失败 %s", err)
	}
	log.Printf("死信交换器 dlx_exchange + 死信队列 dead_letter_queue 已就绪")

	// 声明普通交换器
	err = ch.ExchangeDeclare(
		"normal_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明普通交换器失败 %s", err)
	}

	// 声明普通队列，设置死信交换器和死信路由键
	_, err = ch.QueueDeclare(
		"normal_queue",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    "dlx_exchange",
			"x-dead-letter-routing-key": "dead_letter_routing_key",
		},
	)
	if err != nil {
		log.Fatalf("声明普通队列失败 %s", err)
	}

	// 绑定普通队列到普通交换器
	err = ch.QueueBind(
		"normal_queue",
		"normal_routing_key",
		"normal_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定普通队列失败 %s", err)
	}
	log.Printf("普通交换器 normal_exchange + 普通队列 normal_queue 已就绪（已绑定死信）")

	// 发送消息
	messages := []struct {
		body string
	}{
		{"正常消息-应被消费"},
		{"失败消息-应进死信"},
		{"正常消息-应被消费"},
		{"失败消息-应进死信"},
		{"正常消息-应被消费"},
	}

	for i, m := range messages {
		err = ch.Publish(
			"normal_exchange",
			"normal_routing_key",
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
		fmt.Printf("[%d] 发送消息: %s\n", i+1, m.body)
	}
}
