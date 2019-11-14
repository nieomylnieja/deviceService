package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GivenTickerServiceIsStarted_WhenNewMeasurementComes_ThenReadingIsPassedToSaveChannel(t *testing.T) {
	s := Service{Dao: &Dao{Readings: make(map[int][]DeviceReading), Devices: make(map[int]Device)}}
	s.run()

	input := &RawInput{
		Id:       "0",
		Name:     "TestDevice",
		Interval: "1000",
	}
	devPayload, _ := s.CreateDevicePayload(input)
	m := func(n int) string { return "10C" }
	dev, _ := s.Dao.AddDevice(devPayload)
	s.StartDevice(dev, m)

	time.Sleep(1 * time.Second)
	s.stop()

	expected := "10C"
	result := fmt.Sprintf("%v", s.dao.Readings[0][0].Value)

	assert.Equal(t, expected, result)
}

type mockDao struct {}

func (m *mockDao) AddDevice(device *DevicePayload) (*Device, error) {
	dev := &Device{
		Id:       device.Id,
		Name:     device.Name,
		Value:    "",
		Interval: device.Interval,
		stopChan: nil,
	}
	return dev, nil
}

func Test_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	out := Service{Dao: &mockDao{}}

	device := &DevicePayload{
		Id:       0,
		Name:     "Thermostat",
		Interval: 1000,
	}
	_, err := out.Dao.AddDevice(device)

	assert.NoError(t, err)
}