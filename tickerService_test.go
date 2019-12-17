package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
