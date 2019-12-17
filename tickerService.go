package main

type Publisher interface {
	PublishMeasurement(measurement Measurement, routingKey string)
}

type DevicesTicker struct {
	stopDevices chan bool
}

func NewTickerService() *DevicesTicker {
	return &DevicesTicker{stopDevices: make(chan bool)}
}

func (t *DevicesTicker) Start(allDevices []Device, publisher Publisher) {
	for i := range allDevices {
		go allDevices[i].deviceTicker(publisher, t.stopDevices)
	}
}
