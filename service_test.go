package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GivenTickerServiceIsStarted_WhenNewMeasurementComes_ThenReadingIsPassedToSaveChannel(t *testing.T) {
	dao := &Dao{make(map[int][]DeviceReading), make(map[int]Device)}
	s := DeviceService{dao: dao}
	s.init()
	s.run()

	input := &RawInput{
		Id:       "0",
		Name:     "TestDevice",
		Interval: "1000",
	}
	dev, _ := s.CreateDevicePayload(input)
	m := func(n int) string { return "10C" }
	_ = s.startDevice(dev, m)

	time.Sleep(1 * time.Second)
	s.stop()

	expected := "10C"
	result := fmt.Sprintf("%v", s.dao.Readings[0][0].Value)

	assert.Equal(t, expected, result)
}

type mockDao struct{}

func Test_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	out := DeviceService{dao: mockDao{}}

	device := Device{....}
	result, err := out.addDeviceToDevices(device)

	assert.NoError(t, err)
	assert...
}
