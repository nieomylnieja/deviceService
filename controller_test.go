package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"testing"
)

type MockTickerService struct {
}

func (s *MockTickerService) Start(allDevices []Device, publisher Publisher) {
	return
}

type MockWriterService struct {
}

func (s *MockWriterService) Start(consumer Consumer) {
	return
}

func TestController_AddDevice_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	_, err := c.AddDevice(context.TODO(), &DevicePayload{Name: "test"})

	assert.Equal(t, ErrDao(""), err)
}

func TestController_GetDevice_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	_, err := c.GetDevice(context.TODO(), primitive.NewObjectID().Hex())

	assert.Equal(t, ErrDao(""), err)
}

func TestController_GetPaginatedDevices_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	_, err := c.GetPaginatedDevices(context.TODO(), 0, 2)

	assert.Equal(t, ErrDao(""), err)
}

func TestController_StartTickerService_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	err := c.StartMeasurementService(context.TODO())

	assert.Equal(t, ErrDao(""), err)
}

func TestController_StartMeasurementService_StartsOnlyOnce_CheckedWithError(t *testing.T) {
	c := Controller{
		mainService:   NewService(&mockDao{returnErr: ErrDao("")}),
		amqpService:   &MockAMQPService{},
		tickerService: &MockTickerService{},
		writerService: &MockWriterService{},
		startOnce:     sync.Once{}}
	ctx := context.TODO()

	err := c.StartMeasurementService(ctx)
	assert.Error(t, err)
	err = c.StartMeasurementService(ctx)
	assert.NoError(t, err)
}
