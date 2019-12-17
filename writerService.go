package main

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

type Consumer interface {
	RegisterConsumer() <-chan amqp.Delivery
}

type MeasurementsWriterService struct {
	db           string
	writerClient client.Client
}

func NewWriterService() *MeasurementsWriterService {
	dbAddress := os.Getenv("INFLUXDB_URL")
	dbName := os.Getenv("INFLUXDB_NAME")
	clt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: dbAddress,
	})
	panicOnError(err, "could not initialize influx connection")
	return &MeasurementsWriterService{
		db:           dbName,
		writerClient: clt,
	}
}

func (mws *MeasurementsWriterService) Start(c Consumer) {
	defer mws.closeClient()
	measureChan := c.RegisterConsumer()

	batchPoints, err := mws.batchPointsModel()
	panicOnError(err, "couldn't create batch points model")

	var measurement Measurement
	go func() {
		for msg := range measureChan {
			err = json.Unmarshal(msg.Body, &measurement)
			panicOnError(err, fmt.Sprintf("failed to unmarshall msg: %s", msg.MessageId))
			mws.dbWrite(batchPoints, measurement)
		}
	}()
}

func (mws *MeasurementsWriterService) dbWrite(batchPoints client.BatchPoints, measurement Measurement) {
	point, err := client.NewPoint(
		"deviceValues",
		map[string]string{"deviceId": measurement.Id.String()},
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
