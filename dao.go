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

func (d *Dao) GetAllDevices() ([]Device, error) {
	devices := []Device{}
	for _, dev := range d.data {
		devices = append(devices, dev)
	}
	return devices, nil
}

func (d *Dao) GetPaginatedDevices(limit int, page int) ([]Device, error) {
	lower, upper := setPageBounds(limit, page, len(d.data))
	devices, err := d.GetAllDevices()
	return devices[lower:upper], err
}
