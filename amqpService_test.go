package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMeasurementsAMQP_Connect_ServicePanicsIfNilAddress(t *testing.T) {
	m := MeasurementsAMQP{}

	assert.Panics(t, func() {
		m.connect("")
	})
}
