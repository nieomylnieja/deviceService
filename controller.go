package main

import (
	"context"
	"os"
	"sync"
)

type Controller struct {
	mainService   *Service
	tickerService *TickerService
	writerService *MeasurementsWriterService
	startOnce     sync.Once
}

func NewController(mainService *Service) *Controller {
	return &Controller{
		mainService:   mainService,
		tickerService: NewTickerService(),
		writerService: NewMeasurementsWriterService(os.Getenv("INFLUXDB_URL"),
			os.Getenv("INFLUXDB_NAME")),
		startOnce: sync.Once{},
	}
}

func (c *Controller) StartTickerService(ctx context.Context) error {
	var err error
	c.startOnce.Do(func() {
		err = c.startTickerService(ctx)
		if err != nil {
			return
		}
	})
	return err
}

func (c *Controller) startTickerService(ctx context.Context) error {
	devices, err := c.mainService.GetAllDevices(ctx)
	if err != nil {
		return err
	}

	publish := make(chan Measurement)
	c.tickerService.Start(devices, publish)
	err = c.writerService.Start(publish)

	return err
}

func (c *Controller) GetDevice(id string, ctx context.Context) (*Device, error) {
	return c.mainService.GetDevice(id, ctx)
}

func (c *Controller) AddDevice(devPayload *DevicePayload, ctx context.Context) (*Device, error) {
	return c.mainService.AddDevice(devPayload, ctx)
}

func (c *Controller) GetPaginatedDevices(limit, page int, ctx context.Context) ([]Device, error) {
	return c.mainService.GetPaginatedDevices(limit, page, ctx)
}
