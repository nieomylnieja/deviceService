package main

type TickerService struct {
	stopDevices chan bool
}

func NewTickerService() *TickerService {
	return &TickerService{stopDevices: make(chan bool)}
}

func (t *TickerService) Start(allDevices []Device, publish chan<- Measurement) {
	for i := range allDevices {
		go allDevices[i].deviceTicker(publish, t.stopDevices)
	}
}
