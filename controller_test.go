package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestController_AddDevice(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := NewController(out)

	_, err := c.AddDevice(&DevicePayload{})

	assert.Error(t, ErrDao(""), err)
}

func TestController_GetDevice(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := NewController(out)

	_, err := c.GetDevice(0)

	assert.Error(t, ErrDao(""), err)
}

func TestController_GetPaginatedDevices(t *testing.T) {
	out := NewService(&mockDao{returnErr: ErrDao("")})
	c := NewController(out)

	_, err := c.GetPaginatedDevices(0, 2)

	assert.Error(t, ErrDao(""), err)
}
