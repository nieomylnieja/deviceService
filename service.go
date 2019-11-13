package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"
	"time"
)

type RawInput struct {
	Id       string `json:"id" validate:"gte=0,numeric"`
	Name     string `json:"name" validate:"required,min=2,max=30"`
	Interval string `json:"interval" validate:"required,gt=0,numeric"`
}

type DevicePayload struct {
	Id       int
	Name     string
	Interval int
}

type DeviceDao interface {
	AddDevice(device *DevicePayload) (*Device, error)
}

type DeviceService struct {
	//	Dao             DeviceDao
	dao             *Dao
	DevicesSaveChan chan DeviceInfo
	stopChan        chan bool
}

func (s *DeviceService) init() {
	s.DevicesSaveChan = make(chan DeviceInfo)
	s.stopChan = make(chan bool)
}

func (s *DeviceService) run() {
	go s.tickerService()
}

func (s *DeviceService) stop() {
	for _, dev := range s.dao.Devices {
		s.stopDevice(&dev)
	}
	s.stopChan <- true
}

func (s *DeviceService) tickerService() {
	var deviceInfo DeviceInfo
	var deviceReading DeviceReading
	for {
		select {
		case <-s.stopChan:
			return
		case deviceInfo = <-s.DevicesSaveChan:
			deviceReading = DeviceReading{deviceInfo.Value,
				deviceInfo.When}
			s.dao.Readings[deviceInfo.Id] =
				append(s.dao.Readings[deviceInfo.Id], deviceReading)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (s *DeviceService) parseDeviceInitInput(input *RawInput) (*DevicePayload, error) {
	validate := validator.New()
	validationErrors := validate.Struct(input)
	if validationErrors != nil {
		for _, err := range validationErrors.(validator.ValidationErrors) {
			fmt.Println(err)
		}
		errForwarded := errors.New("input validation failed, device not created")
		return nil, errForwarded
	}
	id, err := strconv.Atoi(input.Id)
	if err != nil {
		return nil, err
	}
	interval, err := strconv.Atoi(input.Interval)
	if err != nil {
		return nil, err
	}
	parsedInput := &DevicePayload{id, input.Name, interval}
	return parsedInput, nil
}

func (s *DeviceService) CreateDevicePayload(input *RawInput) (*DevicePayload, error) {
	devicePayload, err := s.parseDeviceInitInput(input)
	if err != nil {
		return nil, err
	}
	return devicePayload, nil
}

func (s *DeviceService) startDevice(device *Device, getMeasurement measurement) {
	go device.deviceTicker(s, getMeasurement)

}

func (s *DeviceService) AddDevice(device *DevicePayload) (*Device, error) {
	d := &Device{device.Id, device.Name, "",
		device.Interval, make(chan bool)}
	if s.deviceAlreadyExists(d.Id) {
		err := errors.New("The device with given ID already exists!")
		return nil, err
	}
	s.dao.Devices[d.Id] = *d
	return d, nil
}

func (s *DeviceService) UpdateDeviceName(id int, name string) error {
	dev, err := s.GetDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Name = name
	return nil
}

func (s *DeviceService) UpdateDeviceInterval(id int, interval int) error {
	dev, err := s.GetDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Interval = interval
	return nil
}

func (s *DeviceService) updateDeviceValue(d *Device, value string) {
	d.Value = value
}

func (s *DeviceService) stopDevice(d *Device) {
	fmt.Printf("Stopping %s ID:%d...\n", d.Name, d.Id)
	d.stopChan <- true
}

func (s *DeviceService) RemoveDevice(d *Device) {
	s.stopDevice(d)
	delete(s.dao.Devices, d.Id)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%s ID:%d removed.\n", d.Name, d.Id)
}

func (s *DeviceService) GetDevicesList() {
	for _, dev := range s.dao.Devices {
		fmt.Printf("%s -- ID:%d and interval=%d mls\n",
			dev.Name, dev.Id, dev.Interval)
	}
}

func (s *DeviceService) GetDeviceByID(id int) (*Device, error) {
	var err error
	if s.deviceAlreadyExists(id) {
		dev := s.dao.Devices[id]
		return &dev, nil
	}
	err = errors.New("The device with given ID doesn't exist!")
	return nil, err
}

func (s *DeviceService) deviceAlreadyExists(id int) bool {
	for _, dev := range s.dao.Devices {
		if dev.Id == id {
			return true
		}
	}
	return false
}

func (s *DeviceService) GetReadings() []string {
	var fwdReadings []string
	for device, readings := range s.dao.Readings {
		fwdReadings = append(fwdReadings, fmt.Sprintf("Device ID:%d\n", device))
		for _, r := range readings {
			fwdReadings = append(fwdReadings,
				fmt.Sprintf("Nanoseconds: %d -- with value %s\n", r.When.Nanosecond(), r.Value))
		}
	}
	return fwdReadings
}
