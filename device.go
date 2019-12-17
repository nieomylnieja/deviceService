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

func (d *Device) DeviceTicker(p Publisher, stop <-chan bool) {
	routingKey := d.Id.Hex()

	ticker := time.NewTicker(time.Duration(d.Interval) * time.Millisecond)

	for {
		select {
		case <-stop:
			ticker.Stop()
			return
		case <-ticker.C:
			p.PublishMeasurement(
				Measurement{
					Id:    d.Id,
					Value: d.Value,
				}, routingKey)
		}
	}
}
