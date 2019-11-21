package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockDao struct {
	returnValue int
	returnErr   error
	calledTimes int
	device      *Device
}

func (m *mockDao) AddDevice(device *DevicePayload) (int, error) {
	m.calledTimes++
	m.returnValue++
	return m.returnValue, m.returnErr
}

func (m *mockDao) GetDevice(id int) (*Device, error) {
	return m.device, m.returnErr
}

func Test_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	device := &DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	}
	dao := &mockDao{returnValue: 1}
	out := NewService(dao)

	dev, err := out.AddDevice(device)

	assert.NoError(t, err)
	assert.Equal(t, 1, dao.calledTimes)
	assert.NotNil(t, dev)
}

func Test_CorrectDeviceAndDaoFails_ServiceFails(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.AddDevice(&DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	})

	assert.Error(t, err)
}

func Test_GivenIntervalValueBelowZeroOrEqualToZero_ServiceFails(t *testing.T) {
	out := NewService(&mockDao{})

	_, err1 := out.AddDevice(&DevicePayload{Interval: -1})
	_, err2 := out.AddDevice(&DevicePayload{Interval: 0})

	assert.Error(t, err1)
	assert.Error(t, err2)
}

func Test_CorrectPayload_DaoFillsOutId(t *testing.T) {
	out := NewService(&mockDao{})

	dev, err := out.AddDevice(&DevicePayload{Name: "aaa"})

	expected := &Device{Id: 1}

	assert.NoError(t, err)
	assert.Equal(t, expected.Id, dev.Id)
}

func Test_CorrectPayload_ServiceDefaultsInterval(t *testing.T) {
	out := NewService(&mockDao{})

	dev, err := out.AddDevice(&DevicePayload{Name: "aaa"})

	expected := &Device{Interval: 1000}

	assert.NoError(t, err)
	assert.Equal(t, expected.Interval, dev.Interval)
}

func Test_GivenDeviceId_ServiceReturnsDeviceObject(t *testing.T) {
	out := NewService(&mockDao{device: &Device{Name: "name"}})

	dev, err := out.GetDevice(1)

	assert.NoError(t, err)
	assert.Equal(t, &Device{Name: "name"}, dev)
}

func Test_GivenIdThatDoesntExist_ServiceReturnsErrorNotfound(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrNotFound("")})

	_, err := out.GetDevice(1)

	assert.Error(t, err)
}
