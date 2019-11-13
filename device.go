package main

import (
	"fmt"
	"time"
)

type DeviceInfo struct {
	Id    int
	Value string
	When  time.Time
}

type DeviceReading struct {
	Value string
	When  time.Time
}

type Device struct {
	Id       int
	Name     string
	Value    string
	Interval int
	stopChan chan bool
}

type measurement func(int) string

func (d *Device) deviceTicker(s *DeviceService, getMeasurement measurement) {
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
