package main

import (
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewMeasurementsWriterService_ReturnsNewWriterServiceAddress(t *testing.T) {
	actual := NewMeasurementsWriterService(os.Getenv("INFLUX_ADDRESS"))
	err := actual.closeClient()

	expected := &MeasurementsWriterService{
		db: "mydb",
	}

	assert.Equal(t, expected.db, actual.db)
	assert2.NoError(t, err)
}

func TestNewMeasurementsWriterService_GivenWrongAddressServicePanics(t *testing.T) {
	writerService := NewMeasurementsWriterService

	assert2.Panics(t, func() { writerService("abc") })
}
