package main

type TickerService struct{}

func (t TickerService) Start(allDevices []Device, publish chan<- Measurement) {
	stopDevices := make(chan bool)

	for i := range allDevices {
		go allDevices[i].deviceTicker(publish, stopDevices)
	}
}
