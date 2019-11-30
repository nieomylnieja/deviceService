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
	data        []Device
}

func (m *mockDao) AddDevice(device *DevicePayload) (int, error) {
	m.calledTimes++
	m.returnValue++
	return m.returnValue, m.returnErr
}

func (m *mockDao) GetDevice(id int) (*Device, error) {
	return m.device, m.returnErr
}

func (m *mockDao) GetPaginatedDevices(limit int, page int) ([]Device, error) {
	return m.data, m.returnErr
}

func (m *mockDao) GetAllDevices() ([]Device, error) {
	return m.data, m.returnErr
}

func Test_AddDevice_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
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

func Test_AddDevice_CorrectDeviceAndDaoFails_ServiceFails(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.AddDevice(&DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	})

	assert.Error(t, err)
}

func Test_AddDevice_GivenIntervalValueBelowZeroOrEqualToZero_ServiceFails(t *testing.T) {
	out := NewService(&mockDao{})

	_, err1 := out.AddDevice(&DevicePayload{Interval: -1})
	_, err2 := out.AddDevice(&DevicePayload{Interval: 0})

	assert.Error(t, err1)
	assert.Error(t, err2)
}

func Test_Add_Device_CorrectPayload_DaoFillsOutId(t *testing.T) {
	out := NewService(&mockDao{})

	dev, err := out.AddDevice(&DevicePayload{Name: "aaa"})

	expected := &Device{Id: 1}

	assert.NoError(t, err)
	assert.Equal(t, expected.Id, dev.Id)
}

func Test_AddDevice_CorrectPayload_ServiceDefaultsInterval(t *testing.T) {
	out := NewService(&mockDao{})

	dev, err := out.AddDevice(&DevicePayload{Name: "aaa"})

	expected := &Device{Interval: 1000}

	assert.NoError(t, err)
	assert.Equal(t, expected.Interval, dev.Interval)
}

func Test_GetDevice_GivenDaoError_ServiceReturnsErrDao(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetDevice(1)

	assert.Equal(t, ErrDao(""), err)
}

func Test_GetDevice_GivenDeviceId_ServiceReturnsDeviceObject(t *testing.T) {
	out := NewService(&mockDao{device: &Device{Name: "name"}})

	dev, err := out.GetDevice(1)

	assert.NoError(t, err)
	assert.Equal(t, &Device{Name: "name"}, dev)
}

func Test_GetDevice_GivenIdThatDoesntExist_ServiceReturnsNil(t *testing.T) {
	out := NewService(&mockDao{returnErr: nil})

	_, err := out.GetDevice(1)

	assert.NoError(t, err)
}

func Test_GetPaginatedDevices_GivenList_ServiceReturnsList(t *testing.T) {
	out := NewService(&mockDao{data: []Device{{Name: "test name"}}})

	devs, err := out.GetPaginatedDevices(0, 0)

	expected := []Device{{Name: "test name"}}

	assert.NoError(t, err)
	assert.Equal(t, expected, devs)
}

func Test_GetPaginatedDevices_GivenDaoError_ServiceReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetPaginatedDevices(1, 0)

	assert.Equal(t, ErrDao(""), err)
}

func TestService_GetAllDevices_GivenDaoError_ServiceReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetAllDevices()

	assert.Error(t, ErrDao(""), err)
}
