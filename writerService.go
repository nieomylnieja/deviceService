package main

import (
	"fmt"
	influxdb "github.com/influxdata/influxdb-client-go"
	"io"
)

type MeasurementsWriterService struct{}

func (m MeasurementsWriterService) Coonect() error {
	// You can generate a Token from the "Tokens Tab" in the UI
	influx, err := influxdb.New(http: //127.0.0.1:9999, myToken, influxdb.WithHTTPClient(myHTTPClient))
	if err != nil {
		panic(err) // error handling here; normally we wouldn't use fmt but it works for the example
	}
	// Add your app code here
	influx.Close() // closes the client.  After this the client is useless.
}

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
