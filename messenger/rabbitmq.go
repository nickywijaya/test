// Package messenger contains implementation of messenger service.
// Any messaging service should be implemented here.
package messenger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bukalapak/packen/converter"
	"github.com/streadway/amqp"

	gx "github.com/bukalapak/go-xample"
)

// RabbitMQ is a client for rabbitmq service.
// It publishes and receives message via AMQP.
type RabbitMQ struct {
	channel  *amqp.Channel
	messages <-chan amqp.Delivery
	option   RabbitMQOption
}

// RabbitMQOption holds all necessary options for RabbitMQ.
type RabbitMQOption struct {
	Username     string
	Password     string
	Host         string
	VHost        string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	Durable      bool
	Exclusive    bool
}

// NewRabbitMQ returns a pointer of RabbitMQ instance and an error.
// It takes an instance of RabbitMQOption
func NewRabbitMQ(opt RabbitMQOption) (*RabbitMQ, error) {
	// init connection
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", opt.Username, opt.Password, opt.Host, opt.VHost))
	if err != nil {
		return &RabbitMQ{}, err
	}

	// init channel
	channel, err := conn.Channel()
	if err != nil {
		return &RabbitMQ{}, err
	}

	// init exchange declaration
	err = channel.ExchangeDeclare(
		opt.ExchangeName,
		opt.ExchangeType,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return &RabbitMQ{}, err
	}

	// init queue
	queue, err := channel.QueueDeclare(
		opt.RoutingKey,
		opt.Durable,   // durable
		true,          // delete when unused
		opt.Exclusive, // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return &RabbitMQ{}, err
	}

	// set QoS
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return &RabbitMQ{}, err
	}

	// bind the queue
	err = channel.QueueBind(
		queue.Name,
		opt.RoutingKey,
		opt.ExchangeName,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return &RabbitMQ{}, err
	}

	// init consumer
	messages, err := channel.Consume(
		queue.Name,
		"",            // consumer
		false,         // auto-ack
		opt.Exclusive, // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return &RabbitMQ{}, err
	}

	rmq := &RabbitMQ{
		channel:  channel,
		messages: messages,
		option:   opt,
	}
	return rmq, nil
}

// PublishLoginHistory sends LoginHistory to queue
func (r *RabbitMQ) PublishLoginHistory(ctx context.Context, loginHistory gx.LoginHistory) error {
	ctxMap := converter.ContextToMap(ctx)
	ctxMap["Retry"] = int32(ctxMap["Retry"].(int))

	data, err := json.Marshal(loginHistory)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		r.option.ExchangeName,
		r.option.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:     ctxMap,
			ContentType: "application/json",
			Body:        data,
		},
	)

	return err
}

// Listen receives message from queue.
// It takes an pointer of GoXample as its parameter.
// The received message the will be forwarded to GoXample to be saved in database.
func (r *RabbitMQ) Listen(goXample *gx.GoXample) {
	var loginHistory gx.LoginHistory

	for message := range r.messages {
		ctxMap := message.Headers
		ctx, _ := converter.MapToContext(ctxMap)

		err := json.Unmarshal(message.Body, &loginHistory)
		if err != nil {
			continue
		}

		goXample.SaveLoginHistory(ctx, loginHistory)

		message.Ack(false)
	}
}
