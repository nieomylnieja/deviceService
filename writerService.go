package main

import (
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"strconv"
	"time"
)

type MeasurementsWriterService struct {
	db           string
	writerClient client.Client
}

func NewMeasurementsWriterService(dbAddress, dbName string) *MeasurementsWriterService {
	clt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: dbAddress,
	})
	if err != nil {
		log.Panicf("could not initialize influx connection: %s", err.Error())
	}
	return &MeasurementsWriterService{
		db:           dbName,
		writerClient: clt,
	}
}

func (mws *MeasurementsWriterService) Start(publish <-chan Measurement) error {
	defer mws.closeClient()

	batchPoints, err := mws.batchPointsModel()
	if err != nil {
		return err
	}

	go func() {
		for measurement := range publish {
			mws.dbWrite(batchPoints, measurement)
		}
	}()

	return nil
}

func (mws *MeasurementsWriterService) dbWrite(batchPoints client.BatchPoints, measurement Measurement) {
	point, err := client.NewPoint(
		"deviceValues",
		map[string]string{"deviceId": strconv.Itoa(measurement.Id)},
		map[string]interface{}{"value": measurement.Value},
		time.Now())
	if err != nil {
		log.Printf("Could not save %+v: %s", measurement, err.Error())
	}
	batchPoints.AddPoint(point)
	if err = mws.writerClient.Write(batchPoints); err != nil {
		log.Printf("Could not write %+v: %s", measurement, err.Error())
	}
}

func (mws *MeasurementsWriterService) batchPointsModel() (client.BatchPoints, error) {
	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  mws.db,
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return batchPoints, nil
}

func (mws *MeasurementsWriterService) closeClient() error {
	if err := mws.writerClient.Close(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
