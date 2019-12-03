package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	"io"
	"net/http"
	"os"
	"time"
)

type MeasurementsWriterService struct{}

func (m *MeasurementsWriterService) Start(publish <-chan Measurement, w io.Writer) error {
	myClient := http.Client{
		Timeout: 10 * time.Second,
	}

	influx, err := influxdb.New(os.Getenv("INFLUX_ADDRESS"),
		os.Getenv("INFLUX_TOKEN"), influxdb.WithHTTPClient(&myClient))
	if err != nil {
		return err
	}
	defer influx.Close()

	go func() {
		for m := range publish {
			metric := influxdb.NewRowMetric(
				map[string]interface{}{"deviceValues": m.Value},
				"device-metrics",
				map[string]string{"deviceId": fmt.Sprintf("%d", m.Id)},
				time.Now())

			_, err = influx.Write(context.Background(), "DeviceService", "811cf8f341a9dca8", metric)
			if err != nil {
				return
			}
		}
	}()
	return err
}
