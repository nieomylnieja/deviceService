package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DeviceTicker_ChannelReturnsCorrectMeasurement(t *testing.T) {
	mch := make(chan Measurement)
	quit := make(chan bool)
	defer close(mch)
	defer close(quit)

	expected := Measurement{
		Id:    2,
		Value: 24.34,
	}

	d := Device{Id: expected.Id, Value: expected.Value, Interval: 1}

	go d.deviceTicker(mch, quit)
	result := <-mch
	quit <- true

	assert.Equal(t, expected, result)
}
