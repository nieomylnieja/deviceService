package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GivenTickerServiceIsStarted_WhenNewMeasurementComes_ThenReadingIsPassedToSaveChannel(t *testing.T) {
	t.Skip()
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

type mockDao struct {
	returnValue *Device
	returnErr   error
	calledTimes int
}

func (m *mockDao) AddDevice(device *DevicePayload) (*Device, error) {
	m.calledTimes++
	return m.returnValue, m.returnErr
}

func Test_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	device := &DevicePayload{
		Id:       0,
		Name:     "Thermostat",
		Interval: 1000,
	}
	dao := &mockDao{
		returnValue: &Device{
			Id:       device.Id,
			Name:     device.Name,
			Value:    "",
			Interval: device.Interval,
			stopChan: nil,
		},
	}
	out := Service{Dao: dao}

	result, err := out.AddDevice(device)

	assert.NoError(t, err)
	assert.Equal(t, 1, dao.calledTimes)
	assert.NotNil(t, result)
}

func Test_CorrectDeviceAndDaoFails_ServiceFails(t *testing.T) {
	dao := &mockDao{returnErr: fmt.Errorf("test error")}
	out := Service{Dao: dao}

	_, err := out.AddDevice(&DevicePayload{})

	assert.Error(t, err)
}
