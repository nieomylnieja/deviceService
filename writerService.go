package main

import (
	"fmt"
	"io"
)

type MeasurementsWriterService struct{}

func (m MeasurementsWriterService) Start(publish <-chan Measurement, w io.Writer) error {
	var err error
	go func() {
		for m := range publish {
			_, err = fmt.Fprintf(w, "ID:%d -- %f\n", m.Id, m.Value)
			if err != nil {
				return
			}
		}
	}()
	return err
}
