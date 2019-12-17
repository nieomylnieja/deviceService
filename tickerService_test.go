package main

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockAMQPMessage struct {
	body       Measurement
	routingKey string
}

type MockAMQPService struct {
	testChan chan MockAMQPMessage
}

func (mA *MockAMQPService) Start() {
	mA.testChan = make(chan MockAMQPMessage)
}

func (mA *MockAMQPService) PublishMeasurement(measurement Measurement, routingKey string) {
	mA.testChan <- MockAMQPMessage{measurement, routingKey}
}

func (mA *MockAMQPService) RegisterConsumer() <-chan amqp.Delivery {
	return nil
}

func TestTickerService_Start_StopChannelWorksProperly(t *testing.T) {
	ts := NewTickerService()
	mockProducer := &MockAMQPService{}
	devices := []Device{{Name: "test", Interval: 1}, {Name: "test", Interval: 1}}
	defer close(ts.stopDevices)

	ts.Start(devices, mockProducer)
	ts.stopDevices <- true

	assert.Empty(t, ts.stopDevices)
}
