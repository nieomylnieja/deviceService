package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestController_AddDevice_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	_, err := c.AddDevice(&DevicePayload{Name: "test"}, context.TODO())

	assert.Equal(t, ErrDao(""), err)
}

func TestController_GetDevice_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	_, err := c.GetDevice("asd", context.TODO())

	assert.Equal(t, ErrDao(""), err)
}

func TestController_GetPaginatedDevices_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	_, err := c.GetPaginatedDevices(0, 2, context.TODO())

	assert.Equal(t, ErrDao(""), err)
}

func TestController_StartTickerService_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := Controller{mainService: out}

	err := c.StartTickerService(context.TODO())

	assert.Equal(t, ErrDao(""), err)
}
