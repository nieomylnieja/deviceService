package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type MeasurementsAMQP struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	queue      string
}

func NewMeasurementsAMQP() *MeasurementsAMQP {
	m := &MeasurementsAMQP{
		exchange: "measurements",
		queue:    "devices",
	}
	m.connect(os.Getenv("AMQP_URL"))
	return m
}

func (m *MeasurementsAMQP) Start() {
	m.createChannel()
	m.declareExchange()
	m.declareQueue()
	m.bindQueue()
}

func (m *MeasurementsAMQP) RegisterConsumer() <-chan amqp.Delivery {
	delivery, err := m.channel.Consume(
		m.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	panicOnError(err, "failed to register consumer")
	return delivery
}

func (m *MeasurementsAMQP) PublishMeasurement(measurement Measurement, routingKey string) {
	body, err := json.Marshal(measurement)
	panicOnError(err, fmt.Sprintf("couldn't marshall measurement from device: %s", routingKey))

	err = m.channel.Publish(
		m.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	panicOnError(err, fmt.Sprintf("measurement from: %s was not published", routingKey))
}

func (m *MeasurementsAMQP) connect(url string) {
	if url == "" {
		log.Panicln("url must not be empty")
	}
	conn, err := amqp.Dial(url)
	panicOnError(err, fmt.Sprintf("couldn't connect to rabbitMQ server: %s", url))
	m.connection = conn
}

func (m *MeasurementsAMQP) createChannel() {
	ch, err := m.connection.Channel()
	panicOnError(err, "failed to open a channel")

	m.channel = ch
}

func (m *MeasurementsAMQP) declareExchange() {
	err := m.channel.ExchangeDeclare(
		m.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil)
	panicOnError(err, "failed to declare exchange")
}

func (m *MeasurementsAMQP) bindQueue() {
	err := m.channel.QueueBind(
		m.queue,
		"#",
		m.exchange,
		false,
		nil)
	panicOnError(err, "failed to bind queue to exchange")
}

func (m *MeasurementsAMQP) declareQueue() {
	_, err := m.channel.QueueDeclare(
		m.queue,
		true,
		false,
		true,
		false,
		nil,
	)
	panicOnError(err, "failed to declare a queue")
}
