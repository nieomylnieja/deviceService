package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Device struct {
	Id       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `json:"name"`
	Value    float64            `json:"value"`
	Interval int                `json:"interval"`
}

type Measurement struct {
	Id    primitive.ObjectID
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
