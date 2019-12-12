package main

import (
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMeasurementsWriterService_GivenWrongAddressServicePanics(t *testing.T) {
	writerService := NewWriterService

	assert2.Panics(t, func() { writerService("abc", "123") })
}
