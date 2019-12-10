package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type mockDao struct {
	returnValue primitive.ObjectID
	returnErr   error
	calledTimes int
	device      *Device
	data        []Device
}

func (m *mockDao) AddDevice(device *DevicePayload, ctx context.Context) (primitive.ObjectID, error) {
	m.calledTimes++
	return m.returnValue, m.returnErr
}

func (m *mockDao) GetDevice(id string, ctx context.Context) (*Device, error) {
	return m.device, m.returnErr
}

func (m *mockDao) GetPaginatedDevices(limit, page int, ctx context.Context) ([]Device, error) {
	return m.data, m.returnErr
}

func (m *mockDao) GetAllDevices(ctx context.Context) ([]Device, error) {
	return m.data, m.returnErr
}

func Test_AddDevice_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	device := &DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	}
	dao := &mockDao{returnValue: primitive.NewObjectID()}
	out := NewService(dao)

	dev, err := out.AddDevice(device, context.TODO())

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
	}, context.TODO())

	assert.Error(t, err)
}

func Test_AddDevice_GivenIntervalValueBelowZeroOrEqualToZero_ServiceFails(t *testing.T) {
	out := NewService(&mockDao{})

	_, err1 := out.AddDevice(&DevicePayload{Interval: -1}, context.TODO())
	_, err2 := out.AddDevice(&DevicePayload{Interval: 0}, context.TODO())

	assert.Error(t, err1)
	assert.Error(t, err2)
}

func Test_AddDevice_CorrectPayload_ServiceDefaultsInterval(t *testing.T) {
	out := NewService(&mockDao{})

	dev, err := out.AddDevice(&DevicePayload{Name: "aaa"}, context.TODO())

	expected := &Device{Interval: 1000}

	assert.NoError(t, err)
	assert.Equal(t, expected.Interval, dev.Interval)
}

func Test_GetDevice_GivenDaoError_ServiceReturnsErrDao(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetDevice("1", context.TODO())

	assert.Equal(t, ErrDao(""), err)
}

func Test_GetDevice_GivenDeviceId_ServiceReturnsDeviceObject(t *testing.T) {
	id := primitive.NewObjectID()
	out := NewService(&mockDao{device: &Device{Name: "name"}})

	dev, err := out.GetDevice(id.String(), context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, &Device{Name: "name"}, dev)
}

func Test_GetDevice_GivenIdThatDoesntExist_ServiceReturnsNil(t *testing.T) {
	out := NewService(&mockDao{returnErr: nil})

	_, err := out.GetDevice("1", context.TODO())

	assert.NoError(t, err)
}

func Test_GetPaginatedDevices_GivenList_ServiceReturnsList(t *testing.T) {
	out := NewService(&mockDao{data: []Device{{Name: "test name"}}})

	devices, err := out.GetPaginatedDevices(0, 0, context.TODO())

	expected := []Device{{Name: "test name"}}

	assert.NoError(t, err)
	assert.Equal(t, expected, devices)
}

func Test_GetPaginatedDevices_GivenDaoError_ServiceReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetPaginatedDevices(1, 0, context.TODO())

	assert.Equal(t, ErrDao(""), err)
}

func TestService_GetAllDevices_GivenDaoError_ServiceReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetAllDevices(context.TODO())

	assert.Error(t, ErrDao(""), err)
}
