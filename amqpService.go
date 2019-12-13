package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type AMQPService interface {
}

type MeasurementsAMQP struct {
	connection    *amqp.Connection
	measurementCh *amqp.Channel
}

func NewMeasurementsAMQP() *MeasurementsAMQP {
	mR := &MeasurementsAMQP{}
	mR.Start()
	return mR
}

func (m *MeasurementsAMQP) Start() {
	m.connect(os.Getenv("AMQP_URL"))
	m.createChannel()

	exchangeName := "measurements"
	queueName := "publish"

	m.declareExchange(exchangeName)
	m.declareQueue(queueName)
	m.bindQueue(queueName, exchangeName)
}

func (m *MeasurementsAMQP) connect(url string) {
	if url == "" {
		log.Panicln("url most not be empty")
	}
	conn, err := amqp.Dial(url)
	panicOnError(err, fmt.Sprintf("couldn't connect to rabbitMQ server: %s", url))
	m.connection = conn
}

func (m *MeasurementsAMQP) createChannel() {
	ch, err := m.connection.Channel()
	panicOnError(err, "failed to open a channel")

	m.measurementCh = ch
}

func (m *MeasurementsAMQP) declareExchange(name string) {
	err := m.measurementCh.ExchangeDeclare(
		name,
		"direct",
		false,
		false,
		false,
		false,
		nil)
	panicOnError(err, fmt.Sprintf("failed to declare exchange: %s", name))
}

func (m *MeasurementsAMQP) declareQueue(name string) {
	_, err := m.measurementCh.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	panicOnError(err, fmt.Sprintf("failed to declare a queue: %s", name))
}

func (m *MeasurementsAMQP) bindQueue(queueName, exchangeName string) {
	err := m.measurementCh.QueueBind(queueName, "", exchangeName, false, nil)
	panicOnError(err, fmt.Sprintf("failed to bind queue: "+
		"%s to exchange: %s", queueName, exchangeName))
}
