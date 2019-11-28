package main

import "os"

type Controller struct {
	mainService   *Service
	tickerService *TickerService
	writerService *MeasurementsWriterService
}

func NewController(mainService *Service) *Controller {
	return &Controller{
		mainService:   mainService,
		tickerService: NewTickerService(),
		writerService: &MeasurementsWriterService{},
	}
}

func (c *Controller) StartTickerService() error {
	devices, err := c.mainService.GetAllDevices()
	if err != nil {
		return err
	}

	publish := make(chan Measurement)
	c.tickerService.Start(devices, publish)
	err = c.writerService.Start(publish, os.Stdout)

	return err
}

func (c *Controller) GetDevice(id int) (*Device, error) {
	return c.mainService.GetDevice(id)
}

func (c *Controller) AddDevice(devPayload *DevicePayload) (*Device, error) {
	return c.mainService.AddDevice(devPayload)
}

func (c *Controller) GetPaginatedDevices(limit, page int) ([]Device, error) {
	return c.mainService.GetPaginatedDevices(limit, page)
}
