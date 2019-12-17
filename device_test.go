package main

import (
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func Test_DeviceTicker_ChannelReturnsCorrectMeasurement(t *testing.T) {
	mockProducer := &MockAMQPService{}
	mockProducer.Start()
	stop := make(chan bool)
	id := primitive.NewObjectID()
	defer close(stop)

	expected := MockAMQPMessage{
		body: Measurement{
			Id:    id,
			Value: 24.34,
		}, routingKey: id.String(),
	}

	d := Device{Id: expected.body.Id, Value: expected.body.Value, Interval: 1}

	go d.DeviceTicker(mockProducer, stop)
	response := <-mockProducer.testChan
	stop <- true
	close(mockProducer.testChan)
	assert.Equal(t, expected, response)
}
