package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

func main() {
	keys := flag.String("keys", "", "绑定键，逗号分隔（支持通配符，如: kern.*,auth.#,*.critical）")
	flag.Parse()
	if *keys == "" {
		log.Fatal("必须指定 -keys 参数，例如: go run consumer.go -keys 'kern.*,auth.#'")
	}

	bindingKeys := strings.Split(*keys, ",")

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

	queueName := "queue_" + strings.ReplaceAll(
		strings.ReplaceAll(*keys, "*", "all"),
		",", "_",
	)

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败 %s", err)
	}

	for _, key := range bindingKeys {
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
		log.Printf("队列 %s 绑定到 %s (binding key: %s)", q.Name, exChangeName, key)
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

	fmt.Printf("[%s] 准备接收消息\n", *keys)
	for d := range msgs {
		log.Printf("[%s] 收到消息 [%s]: %s\n", *keys, d.RoutingKey, d.Body)
	}
}
