package main

import (
	"fmt"
	"log"
	"strconv"

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

	err = ch.ExchangeDeclare(
		"rpc_request_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明请求交换器失败 %s", err)
	}

	_, err = ch.QueueDeclare(
		"rpc_request_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明请求队列失败 %s", err)
	}

	err = ch.QueueBind(
		"rpc_request_queue",
		"rpc_request",
		"rpc_request_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定请求队列失败 %s", err)
	}

	msgs, err := ch.Consume(
		"rpc_request_queue",
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

	fmt.Println("消费者启动，等待请求...")

	for d := range msgs {
		n, err := strconv.Atoi(string(d.Body))
		if err != nil {
			d.Nack(false, false)
			continue
		}

		result := n * n
		log.Printf("[处理] n=%d, 结果=%d, 回复到 %s\n", n, result, d.ReplyTo)

		err = ch.Publish(
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(strconv.Itoa(result)),
			},
		)
		if err != nil {
			log.Printf("回复失败: %v", err)
			d.Nack(false, false)
			continue
		}

		d.Ack(false)
	}
}
