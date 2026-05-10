package main

import (
	"fmt"
	"log"
	"math/rand"
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

	// 声明请求交换器
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

	// 声明请求队列
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

	// 声明回调队列（临时，连接断开自动删除）
	callbackQueue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明回调队列失败 %s", err)
	}
	log.Printf("回调队列已创建: %s", callbackQueue.Name)

	// 监听回调队列
	callbackMsgs, err := ch.Consume(
		callbackQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("注册回调消费者失败 %s", err)
	}

	// 发送请求消息，设置 ReplyTo 和 CorrelationId
	numbers := []int{10, 20, 30, 40, 50}
	for _, n := range numbers {
		corrId := strconv.Itoa(rand.Intn(100000))

		err = ch.Publish(
			"rpc_request_exchange",
			"rpc_request",
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				ReplyTo:       callbackQueue.Name,
				CorrelationId: corrId,
				Body:          []byte(strconv.Itoa(n)),
			},
		)
		if err != nil {
			log.Printf("发送失败: %v", err)
			return
		}
		fmt.Printf("[请求] n=%d, corrId=%s → 等待回调\n", n, corrId)
	}

	// 等待回调响应
	received := 0
	for d := range callbackMsgs {
		log.Printf("[回调] corrId=%s, 结果=%s\n", d.CorrelationId, d.Body)
		received++
		if received == len(numbers) {
			break
		}
	}
}
