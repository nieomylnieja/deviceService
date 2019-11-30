package main

import (
	"time"
)

type Device struct {
	Id       int
	Name     string
	Value    float64
	Interval int
}

type Measurement struct {
	Id    int
	Value float64
}

func (d *Device) deviceTicker(publish chan<- Measurement, stop <-chan bool) {
	ticker := time.NewTicker(time.Duration(d.Interval) * time.Millisecond)

	for {
		select {
		case <-stop:
			ticker.Stop()
			return
		case <-ticker.C:
			publish <- Measurement{
				Id:    d.Id,
				Value: d.Value,
			}
		}
	}
}
