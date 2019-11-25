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

func (d *Device) deviceTicker(mch chan Measurement, quit chan bool) {
	ticker := time.NewTicker(time.Duration(d.Interval) * time.Millisecond)

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			mch <- Measurement{
				Id:    d.Id,
				Value: d.Value,
			}
		}
	}
}
