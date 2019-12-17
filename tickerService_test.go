package main

import (
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

func TestTickerService_Start_StopChannelWorksProperly(t *testing.T) {
	ts := NewTickerService()
	mockProducer := &MockAMQPService{}
	devices := []Device{{Name: "test", Interval: 1}, {Name: "test", Interval: 1}}
	publish := make(chan Measurement)
	defer close(publish)
	defer close(ts.stopDevices)

	ts.Start(devices, mockProducer)
	ts.stopDevices <- true

	assert.Empty(t, publish)
	assert.Empty(t, ts.stopDevices)
}
