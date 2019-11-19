package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockDao struct {
	returnValue int
	returnErr   error
	calledTimes int
}

func (m *mockDao) AddDevice(device *DevicePayload) (int, error) {
	m.calledTimes++
	m.returnValue++
	return m.returnValue, m.returnErr
}

func Test_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	device := &DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	}
	dao := &mockDao{returnValue: 1}

	out := Service{Dao: dao}

	result, err := out.AddDevice(device)

	assert.NoError(t, err)
	assert.Equal(t, 1, dao.calledTimes)
	assert.NotNil(t, result)
}

func Test_CorrectDeviceAndDaoFails_ServiceFails(t *testing.T) {
	dao := &mockDao{returnErr: fmt.Errorf("test error")}
	out := Service{Dao: dao}

	_, err := out.AddDevice(&DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	})

	assert.Error(t, err)
}

func Test_GivenIntervalValueBelowZeroOrEqualToZero_ServiceFails(t *testing.T) {
	out := Service{Dao: &mockDao{}}

	_, err1 := out.AddDevice(&DevicePayload{Interval: -1})
	_, err2 := out.AddDevice(&DevicePayload{Interval: 0})
	assert.Error(t, err1)
	assert.Error(t, err2)
}

func Test_CorrectPayload_DaoFillsOutId_ServiceDefaultsInterval(t *testing.T) {
	dao := &mockDao{}
	out := Service{Dao: dao}

	_, err := out.AddDevice(&DevicePayload{Name: "aaa"})

	expected := &Device{
		Interval: 1000,
		Name:     "aaa",
		Id:       1,
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, out.tempDevice)
}
