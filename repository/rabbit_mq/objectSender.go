package rabbit_mq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQSender struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

type TaskMessageRMQ struct {
	TaskID     string `json:"task_id"`
	Translator string `json:"translator"`
	Code       string `json:"code"`
}

func NewRabbitMQSender(amqpURL, queueName string) (*RabbitMQSender, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to rabbitMQ:  %s", err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &RabbitMQSender{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (sender *RabbitMQSender) Send(task TaskMessageRMQ) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = sender.channel.Publish("", sender.queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (sender *RabbitMQSender) Close() error {
	err := sender.channel.Close()
	if err != nil {
		return err
	}
	err = sender.connection.Close()
	return err
}
