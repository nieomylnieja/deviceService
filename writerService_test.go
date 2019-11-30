package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMeasurementsWriterService_Start_OutputsCorrectStringAndCloses(t *testing.T) {
	ws := MeasurementsWriterService{}
	publish := make(chan Measurement)
	measurement := Measurement{
		Id:    1,
		Value: 2,
	}
	var buf bytes.Buffer

	err := ws.Start(publish, &buf)
	publish <- measurement
	close(publish)
	expected := fmt.Sprintf("ID:%d -- %f\n", measurement.Id, measurement.Value)

	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}
