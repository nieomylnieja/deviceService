package main

import (
	"github.com/fatih/structs"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"os"
	"strconv"
	"time"
)

type MeasurementsWriterService struct {
	mydb         string
	writerClient client.Client
}

func NewMeasurementsWriterService() *MeasurementsWriterService {
	clt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: os.Getenv("INFLUX_ADDRESS"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return &MeasurementsWriterService{
		mydb:         "mydb",
		writerClient: clt,
	}
}

func (mws *MeasurementsWriterService) insert(batchPoints client.BatchPoints, measurement Measurement) {
	defer mws.writerClient.Close()

	point, err := client.NewPoint(
		"deviceValues",
		map[string]string{"deviceId": strconv.Itoa(measurement.Id)},
		structs.Map(measurement.Value),
		time.Now())
	if err != nil {
		log.Println(err)
	}
	batchPoints.AddPoint(point)
	if err = mws.writerClient.Write(batchPoints); err != nil {
		log.Println(err)
	}
	if err = mws.writerClient.Close(); err != nil {
		log.Println(err)
	}
}

func (mws *MeasurementsWriterService) Start(publish <-chan Measurement) error {
	var err error

	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  mws.mydb,
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
	}

	go func() {
		for measurement := range publish {
			mws.insert(batchPoints, measurement)
		}
	}()
	return err
}
