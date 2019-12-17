package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type MockAMQPService struct {
}

func (mA *MockAMQPService) Start() {
	return
}

func (mA *MockAMQPService) PublishMeasurement(measurement Measurement, routingKey string) {
	return
}

func Test_DeviceTicker_ChannelReturnsCorrectMeasurement(t *testing.T) {
	mockProducer := &MockAMQPService{}
	stop := make(chan bool)
	id := primitive.NewObjectID()
	defer close(stop)

	expected := Measurement{
		Id:    id,
		Value: 24.34,
	}

	d := Device{Id: expected.Id, Value: expected.Value, Interval: 1}

	go d.deviceTicker(mockProducer, stop)
	stop <- true
}
