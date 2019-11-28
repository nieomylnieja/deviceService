package main

import (
	"fmt"
	"sync"
)

type MeasurementsWriterService struct {
	once sync.Once
}

type WriterService interface {
	Start(publish <-chan Measurement)
}

func (m MeasurementsWriterService) Start(publish <-chan Measurement) {
	go func() {
		for m := range publish {
			fmt.Printf("ID:%d -- %f\n", m.Id, m.Value)
		}
	}()
}
