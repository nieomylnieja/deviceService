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

func (m *mockDao) AddDevice(ctx context.Context, device *DevicePayload) (primitive.ObjectID, error) {
	m.calledTimes++
	return m.returnValue, m.returnErr
}

func (m *mockDao) GetDevice(ctx context.Context, id primitive.ObjectID) (*Device, error) {
	return m.device, m.returnErr
}

func (m *mockDao) GetPaginatedDevices(ctx context.Context, limit, page int) ([]Device, error) {
	return m.data, m.returnErr
}

func (m *mockDao) GetAllDevices(ctx context.Context) ([]Device, error) {
	return m.data, m.returnErr
}

func TestService_AddDevice_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	device := &DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	}
	dao := &mockDao{returnValue: primitive.NewObjectID()}
	out := NewService(dao)

	dev, err := out.AddDevice(context.TODO(), device)

	assert.NoError(t, err)
	assert.Equal(t, 1, dao.calledTimes)
	assert.NotNil(t, dev)
}

func TestService_AddDevice_CorrectDeviceAndDaoFails_ServiceFails(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.AddDevice(context.TODO(),
		&DevicePayload{
			Value:    10.23,
			Name:     "Thermostat",
			Interval: 1000,
		})

	assert.Error(t, err)
}

func TestService_AddDevice_GivenIntervalValueBelowZeroOrEqualToZero_ServiceFails(t *testing.T) {
	out := NewService(&mockDao{})

	_, err1 := out.AddDevice(context.TODO(), &DevicePayload{Interval: -1})
	_, err2 := out.AddDevice(context.TODO(), &DevicePayload{Interval: 1})

	assert.Error(t, err1)
	assert.Error(t, err2)
}

func TestService_AddDevice_CorrectPayload_ServiceDefaultsInterval(t *testing.T) {
	out := NewService(&mockDao{})

	dev, err := out.AddDevice(context.TODO(), &DevicePayload{Name: "aaa"})

	expected := &Device{Interval: 1000}

	assert.NoError(t, err)
	assert.Equal(t, expected.Interval, dev.Interval)
}

func TestService_GetDevice_GivenDaoError_ServiceReturnsErrDao(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	id := primitive.NewObjectID().Hex()

	_, err := out.GetDevice(context.TODO(), id)

	assert.Equal(t, ErrDao(""), err)
}

func TestService_GetDevice_GivenDeviceId_ServiceReturnsDeviceObject(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	out := NewService(&mockDao{device: &Device{Name: "name"}})

	dev, err := out.GetDevice(context.TODO(), id)

	assert.NoError(t, err)
	assert.Equal(t, &Device{Name: "name"}, dev)
}

func TestService_GetDevice_GivenIdThatDoesntExist_ServiceReturnsNil(t *testing.T) {
	out := NewService(&mockDao{returnErr: nil})
	id := primitive.NewObjectID().Hex()

	_, err := out.GetDevice(context.TODO(), id)

	assert.NoError(t, err)
}

func TestService_GetPaginatedDevices_GivenList_ServiceReturnsList(t *testing.T) {
	out := NewService(&mockDao{data: []Device{{Name: "test name"}}})

	devices, err := out.GetPaginatedDevices(context.TODO(), 0, 0)

	expected := []Device{{Name: "test name"}}

	assert.NoError(t, err)
	assert.Equal(t, expected, devices)
}

func TestService_GetPaginatedDevices_GivenDaoError_ServiceReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetPaginatedDevices(context.TODO(), 1, 0)

	assert.Equal(t, ErrDao(""), err)
}

func TestService_GetAllDevices_GivenDaoError_ServiceReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})

	_, err := out.GetAllDevices(context.TODO())

	assert.Error(t, ErrDao(""), err)
}
