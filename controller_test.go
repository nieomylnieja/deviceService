package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestController_AddDevice_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := NewController(out)

	_, err := c.AddDevice(&DevicePayload{Name: "test"})

	assert.Equal(t, ErrDao(""), err)
}

func TestController_GetDevice_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := NewController(out)

	_, err := c.GetDevice(0)

	assert.Equal(t, ErrDao(""), err)
}

func TestController_GetPaginatedDevices_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := NewController(out)

	_, err := c.GetPaginatedDevices(0, 2)

	assert.Equal(t, ErrDao(""), err)
}

func TestController_StartTickerService_GivenDaoError_ControllerReturnsError(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := NewController(out)

	err := c.StartTickerService()

	assert.Equal(t, ErrDao(""), err)
}
