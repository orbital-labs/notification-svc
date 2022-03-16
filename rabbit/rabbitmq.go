// Package messenger contains implementation of messenger service.
// Any messaging service should be implemented here.
package rabbit

import (
	"fmt"
	"log"

	"github.com/assembla/cony"
	"github.com/streadway/amqp"
)

// RabbitMQ is a client for rabbitmq service.
// It publishes and receives message via AMQP.
type RabbitMQ struct {
	option    RabbitMQOption
	client    *cony.Client
	publisher *cony.Publisher
	consumer  *cony.Consumer
}

// RabbitMQOption holds all necessary options for RabbitMQ.
type RabbitMQOption struct {
	RabbitMQURL  string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	Durable      bool
	Exclusive    bool
}

func NewRabbitMQPublisher(opt RabbitMQOption) *RabbitMQ {
	url := opt.RabbitMQURL
	cli := cony.NewClient(
		cony.URL(url),
		cony.Backoff(cony.DefaultBackoff),
	)

	exc := cony.Exchange{
		Name:       opt.ExchangeName,
		Kind:       opt.ExchangeType,
		Durable:    true,
		AutoDelete: false,
	}

	cli.Declare([]cony.Declaration{
		cony.DeclareExchange(exc),
	})

	publisher := cony.NewPublisher(opt.ExchangeName, opt.RoutingKey)
	cli.Publish(publisher)

	rmq := &RabbitMQ{
		option:    opt,
		client:    cli,
		publisher: publisher,
	}

	return rmq
}

func NewRabbitMQConsumer(opt RabbitMQOption) *RabbitMQ {
	url := opt.RabbitMQURL
	cli := cony.NewClient(
		cony.URL(url),
		cony.Backoff(cony.DefaultBackoff),
	)

	exc := cony.Exchange{
		Name:       opt.ExchangeName,
		Kind:       opt.ExchangeType,
		Durable:    true,
		AutoDelete: false,
	}

	que := &cony.Queue{
		AutoDelete: false,
		Name:       opt.RoutingKey,
	}

	bnd := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      opt.RoutingKey,
	}

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd),
	})

	cns := cony.NewConsumer(
		que,
		cony.AutoAck(),
	)
	cli.Consume(cns)

	rmq := &RabbitMQ{
		option:   opt,
		client:   cli,
		consumer: cns,
	}

	return rmq
}

func (r *RabbitMQ) ServePublish() {
	for r.client.Loop() {
		select {
		case err := <-r.client.Errors():
			fmt.Printf("Consumer error: %v\n", err)
		}
	}
}

func (r *RabbitMQ) ServeConsume(f func(amqp.Delivery)) {
	for r.client.Loop() {
		select {
		case msg := <-r.consumer.Deliveries():
			log.Printf("Received body: %q\n", msg.Body)
			f(msg)
		case err := <-r.consumer.Errors():
			fmt.Printf("Consumer error: %v\n", err)
		case err := <-r.client.Errors():
			fmt.Printf("Client error: %v\n", err)
		}
	}
}

// PublishImportRequest sends ImportRequest to queue
func (r *RabbitMQ) Publish(data []byte) error {
	err := r.publisher.Publish(
		amqp.Publishing{
			Headers:      nil,
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
			Priority:     1,
		},
	)
	return err
}

func (r *RabbitMQ) PublishUnprioritized(data []byte) error {
	err := r.publisher.Publish(
		amqp.Publishing{
			Headers:      nil,
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
			Priority:     0,
		},
	)
	return err
}

func (r *RabbitMQ) ListenRequest(message amqp.Delivery) {
	fmt.Print("Pesan Diterima")
}
