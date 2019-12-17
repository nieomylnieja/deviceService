package main

import (
	"context"
	"sync"
)

type MainService interface {
	AddDevice(ctx context.Context, payload *DevicePayload) (*Device, error)
	GetDevice(ctx context.Context, id string) (*Device, error)
	GetPaginatedDevices(ctx context.Context, limit, page int) ([]Device, error)
	GetAllDevices(ctx context.Context) ([]Device, error)
}

type TickerService interface {
	Start(allDevices []Device, publisher Publisher)
}

type WriterService interface {
	Start(consumer Consumer)
}

type AMQPService interface {
	Start()
	Publisher
	Consumer
}

type Controller struct {
	mainService   MainService
	amqpService   AMQPService
	tickerService TickerService
	writerService WriterService
	startOnce     sync.Once
}

func NewController(
	mainService MainService,
	amqpService AMQPService,
	writerService WriterService,
	tickerService TickerService) *Controller {

	return &Controller{
		mainService:   mainService,
		amqpService:   amqpService,
		writerService: writerService,
		tickerService: tickerService,
		startOnce:     sync.Once{},
	}
}

func (c *Controller) StartMeasurementService(ctx context.Context) error {
	var err error
	c.startOnce.Do(func() {
		err = c.startMeasurementServices(ctx)
		if err != nil {
			return
		}
	})
	return err
}

func (c *Controller) startMeasurementServices(ctx context.Context) error {
	devices, err := c.mainService.GetAllDevices(ctx)
	if err != nil {
		return err
	}

	c.amqpService.Start()
	c.tickerService.Start(devices, c.amqpService)
	c.writerService.Start(c.amqpService)

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
