package messenger

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	option     RabbitMQOption
}

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

func NewRabbitMQPublisher(opt RabbitMQOption) (*RabbitMQPublisher, error) {
	// init connection
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", opt.Username, opt.Password, opt.Host, opt.VHost))
	if err != nil {
		return &RabbitMQPublisher{}, err
	}

	// init channel
	channel, err := conn.Channel()
	if err != nil {
		return &RabbitMQPublisher{}, err
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
		return &RabbitMQPublisher{}, err
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
		return &RabbitMQPublisher{}, err
	}

	// set QoS
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return &RabbitMQPublisher{}, err
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
		return &RabbitMQPublisher{}, err
	}

	pub := &RabbitMQPublisher{
		Connection: conn,
		Channel:    channel,
		Queue:      queue,
		option:     opt,
	}
	return pub, nil
}

func (r *RabbitMQPublisher) Publish(contentType string, data []byte) error {
	err := r.Channel.Publish(
		r.option.ExchangeName,
		r.option.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        data,
		},
	)

	return err
}
