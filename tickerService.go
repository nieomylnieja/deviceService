package main

type TickerService interface {
	Start(allDevices []Device, publish chan<- Measurement)
}

type DevicesTicker struct {
	stopDevices chan bool
}

func NewTickerService() *DevicesTicker {
	return &DevicesTicker{stopDevices: make(chan bool)}
}

func (t *DevicesTicker) Start(allDevices []Device, publish chan<- Measurement) {
	for i := range allDevices {
		go allDevices[i].deviceTicker(publish, t.stopDevices)
	}
}
