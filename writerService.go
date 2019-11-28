package main

import (
	"fmt"
)

type MeasurementsWriterService struct{}

func (m MeasurementsWriterService) Start(publish <-chan Measurement) {
	go func() {
		for m := range publish {
			fmt.Printf("ID:%d -- %f\n", m.Id, m.Value)
		}
	}()
}
