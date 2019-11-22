package main

type Dao struct {
	data map[int]Device

	indexer int
}

func NewDao() *Dao {
	d := Dao{data: make(map[int]Device)}
	return &d
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

func (d *Dao) GetDevice(id int) (*Device, error) {
	if device, ok := d.data[id]; ok {
		return &device, nil
	}
	return nil, nil
}

func (d *Dao) GetManyDevices(limit int, page int) ([]Device, error) {
	if limit == 0 {
		allDevices := make([]Device, d.indexer)
		for _, device := range d.data {
			allDevices = append(allDevices, device)
		}
		return allDevices, nil
	}

	if limit*page > d.indexer {
		return []Device{}, nil
	}

	if page*limit+limit >= d.indexer {
		someDevices := make([]Device, d.indexer%limit)
		for _, device := range d.data {
			someDevices = append(someDevices, device)
		}
		return someDevices, nil
	}

	someDevices := make([]Device, limit)
	for _, device := range d.data {
		someDevices = append(someDevices, device)
	}
	return someDevices, nil
}
