package main

import (
	"fmt"
	"time"
)

type DeviceInfo struct {
	Id    int
	Value float64
	When  time.Time
}

type DeviceReading struct {
	Value float64
	When  time.Time
}

type Device struct {
	Id       int
	Name     string
	Value    float64
	Interval int
	stopChan chan bool
}

type measurement func(int) float64

func (d *Device) deviceTicker(s *Service, getMeasurement measurement) {
	ticker := time.NewTicker(time.Duration(d.Interval) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(d.Interval) * time.Millisecond)
			s.updateDeviceValue(d, getMeasurement(d.Id))
			s.DevicesSaveChan <- DeviceInfo{d.Id, d.Value, time.Now()}
		case <-d.stopChan:
			ticker.Stop()
			fmt.Printf("...%s ID:%d stopped!\n", d.Name, d.Id)
			return
		}
	}
}
