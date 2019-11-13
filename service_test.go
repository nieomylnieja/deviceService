package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GivenTickerServiceIsStarted_WhenNewMeasurementComes_ThenReadingIsPassedToSaveChannel(t *testing.T) {
	dao := &DataAccessObject{make(map[int][]DeviceReading)}
	s := Service{}
	s.init(dao)
	s.run()

	input := &RawInput{
		Id:       "0",
		Name:     "TestDevice",
		Interval: "1000",
	}
	dev, _ := s.createDevice(input)
	m := func(n int) string { return "10C"}
	_ = s.startDevice(dev, m)

	time.Sleep(1 * time.Second)
	s.stop()

	expected := "10C"
	result := fmt.Sprintf("%v", s.DevicesReadings.data[0][0].Value)
	fmt.Println(expected)
	assert.Equal(t, expected, result)
}