package main

type Dao struct {
	data map[int]Device

	indexer int
}

func (d *Dao) AddDevice(device *DevicePayload) (int, error) {
	d.indexer++
	dev := &Device{
		Id:       d.indexer,
		Name:     device.Name,
		Value:    device.Value,
		Interval: device.Interval,
	}

	d.data[dev.Id] = *dev
	return dev.Id, nil
}
