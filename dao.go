package main

type Dao struct {
	Readings map[int][]DeviceReading
	Devices  map[int]Device

	indexer int
}

func (d *Dao) AddDevice(device *DevicePayload) (int, error) {
	d.indexer++
	dev := &Device{
		Id:       d.indexer,
		Name:     device.Name,
		Value:    device.Value,
		Interval: device.Interval,
		stopChan: make(chan bool),
	}

	d.Devices[dev.Id] = *dev
	return dev.Id, nil
}
