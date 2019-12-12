package main

import (
	"context"
	"sync"
)

type Controller struct {
	mainService   *Service
	tickerService TickerService
	writerService WriterService
	startOnce     sync.Once
}

func NewController(mainService *Service, writerService WriterService,
	tickerService TickerService) *Controller {
	return &Controller{
		mainService:   mainService,
		writerService: writerService,
		tickerService: tickerService,
		startOnce:     sync.Once{},
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

func (c *Controller) GetDevice(ctx context.Context, id string) (*Device, error) {
	return c.mainService.GetDevice(ctx, id)
}

func (c *Controller) AddDevice(ctx context.Context, devPayload *DevicePayload) (*Device, error) {
	return c.mainService.AddDevice(ctx, devPayload)
}

func (c *Controller) GetPaginatedDevices(ctx context.Context, limit, page int) ([]Device, error) {
	return c.mainService.GetPaginatedDevices(ctx, limit, page)
}
