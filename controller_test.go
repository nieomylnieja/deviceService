package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

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
