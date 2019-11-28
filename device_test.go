package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DeviceTicker_ChannelReturnsCorrectMeasurement(t *testing.T) {
	publish := make(chan Measurement)
	stop := make(chan bool)
	defer close(publish)
	defer close(stop)

	expected := Measurement{
		Id:    2,
		Value: 24.34,
	}

	d := Device{Id: expected.Id, Value: expected.Value, Interval: 1}

	go d.deviceTicker(publish, stop)
	result := <-publish
	stop <- true

	assert.Equal(t, expected, result)
}
