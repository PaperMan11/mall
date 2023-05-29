package rabbitmq

import (
	"encoding/json"
	"fmt"
	"mall/model"

	"github.com/streadway/amqp"
)

type SkAmqp interface {
	SendToQueue(qName string, req *model.SkillReq2MQ) error
	ReceiveFromQueue(qName string) (*amqp.Channel, <-chan amqp.Delivery, error)
}

type customSkAmqp struct {
	amqpConn *amqp.Connection
}

func NewSkAmqp(amqpConn *amqp.Connection) SkAmqp {
	return &customSkAmqp{
		amqpConn: amqpConn,
	}
}

func (c *customSkAmqp) SendToQueue(qName string, req *model.SkillReq2MQ) error {
	data, _ := json.Marshal(req)
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return failOnError(err, "Failed to open a channel")
	}
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		qName, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		return failOnError(err, "Failed to declare a queue")
	}
	// 将消息发布到声明的队列中
	if err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immmediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	); err != nil {
		failOnError(err, "Failed to publish a message")
	}
	return nil
}

// 单独的任务
func (c *customSkAmqp) ReceiveFromQueue(qName string) (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, nil, failOnError(err, "Failed to open a channel")
	}

	// 声明队列
	q, err := ch.QueueDeclare(
		qName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, failOnError(err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // queue
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args （设置参数）可以用于关联 不同的 交换机（例如：无法消费的任务转发给->死信队列）
	)
	if err != nil {
		return nil, nil, failOnError(err, "Failed to register a consumer")
	}

	return ch, msgs, nil
}

func failOnError(err error, msg string) error {
	return fmt.Errorf("%s: %s", msg, err)
}
