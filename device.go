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

func (d *Device) deviceTicker(mch chan Measurement, stop chan bool) {
	ticker := time.NewTicker(time.Duration(d.Interval) * time.Millisecond)

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			mch <- Measurement{
				Id:    d.Id,
				Value: d.Value,
			}
		}
	}
}
