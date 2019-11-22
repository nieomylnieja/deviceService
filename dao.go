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
	i := 0
	if limit == 0 {
		devices := make([]Device, 0, d.indexer)
		for _, dev := range d.data {
			if i >= d.indexer {
				break
			}
			devices = append(devices, dev)
			i++
		}
		return devices, nil
	}

	if limit*page > d.indexer {
		return []Device{}, nil
	}

	if page*limit+limit >= d.indexer {
		devices := make([]Device, 0, d.indexer%limit)
		for _, dev := range d.data {
			if i >= d.indexer%limit {
				break
			}
			devices = append(devices, dev)
			i++
		}
		return devices, nil
	}

	devices := make([]Device, 0, limit)
	for _, dev := range d.data {
		if i >= limit {
			break
		}
		devices = append(devices, dev)
		i++
	}
	return devices, nil
}
