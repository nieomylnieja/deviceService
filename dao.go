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
	for _, device := range d.data {
		if device.Id == id {
			return &device, nil
		}
	}
	return nil, nil
}
