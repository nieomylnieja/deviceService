package main

import "errors"

type Dao struct {
	Readings map[int][]DeviceReading
	Devices  map[int]Device
}

func (d *Dao) AddDevice(device *DevicePayload) (*Device, error) {
	dev := &Device{device.Id, device.Name, "",
		device.Interval, make(chan bool)}
	if d.deviceAlreadyExists(dev.Id) {
		err := errors.New("The device with given ID already exists!")
		return nil, err
	}
	d.Devices[dev.Id] = *dev
	return dev, nil
}

func (d *Dao) deviceAlreadyExists(id int) bool {
	for _, dev := range d.Devices {
		if dev.Id == id {
			return true
		}
	}
	return false
}