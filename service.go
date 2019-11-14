package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type RawInput struct {
	Id       string `json:"id" validate:"gte=0,numeric"`
	Name     string `json:"name" validate:"required,min=2,max=30"`
	Interval string `json:"interval" validate:"required,gt=0,numeric"`
}

type DevicePayload struct {
	Id       int    `json:"id" validate:"gte=0,numeric"`
	Name     string `json:"name" validate:"required,min=2,max=30"`
	Interval int    `json:"interval" validate:"required,gt=0,numeric"`
}

type DeviceDao interface {
	AddDevice(device *DevicePayload) (*Device, error)
}

type Service struct {
	Dao             DeviceDao
	dao             *Dao
	DevicesSaveChan chan DeviceInfo
	stopChan        chan bool
}

func (s *Service) run() {
	s.dao = &Dao{
		Readings: make(map[int][]DeviceReading),
		Devices:  make(map[int]Device),
	}
	s.DevicesSaveChan = make(chan DeviceInfo)
	s.stopChan = make(chan bool)
	go s.tickerService()
}

func (s *Service) stop() {
	for _, dev := range s.dao.Devices {
		s.stopDevice(&dev)
	}
	s.stopChan <- true
}

func (s *Service) tickerService() {
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

func (s *Service) parseDeviceInitInput(input *RawInput) (*DevicePayload, error) {
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

func (s *Service) CreateDevicePayload(input *RawInput) (*DevicePayload, error) {
	devicePayload, err := s.parseDeviceInitInput(input)
	if err != nil {
		return nil, err
	}
	return devicePayload, nil
}

func (s *Service) StartDevice(device *Device, getMeasurement measurement) {
	go device.deviceTicker(s, getMeasurement)

}

func (s *Service) UpdateDeviceName(id int, name string) error {
	dev, err := s.GetDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Name = name
	return nil
}

func (s *Service) UpdateDeviceInterval(id int, interval int) error {
	dev, err := s.GetDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Interval = interval
	return nil
}

func (s *Service) updateDeviceValue(d *Device, value string) {
	d.Value = value
}

func (s *Service) stopDevice(d *Device) {
	fmt.Printf("Stopping %s ID:%d...\n", d.Name, d.Id)
	d.stopChan <- true
}

func (s *Service) RemoveDevice(d *Device) {
	s.stopDevice(d)
	delete(s.dao.Devices, d.Id)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%s ID:%d removed.\n", d.Name, d.Id)
}

func (s *Service) GetDevicesList() {
	for _, dev := range s.dao.Devices {
		fmt.Printf("%s -- ID:%d and interval=%d mls\n",
			dev.Name, dev.Id, dev.Interval)
	}
}

func (s *Service) GetDeviceByID(id int) (*Device, error) {
	var err error
	if s.deviceAlreadyExists(id) {
		dev := s.dao.Devices[id]
		return &dev, nil
	}
	err = errors.New("The device with given ID doesn't exist!")
	return nil, err
}

func (s *Service) deviceAlreadyExists(id int) bool {
	for _, dev := range s.dao.Devices {
		if dev.Id == id {
			return true
		}
	}
	return false
}

func (s *Service) GetReadings() []string {
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

func (s *Service) validate(payload *DevicePayload) error {
	validate := validator.New()
	validationErrors := validate.Struct(payload)
	if validationErrors != nil {
		for _, err := range validationErrors.(validator.ValidationErrors) {
			fmt.Println(err)
		}
		errForwarded := errors.New("input validation failed, device not created")
		return errForwarded
	}
	return nil
}

func (s *Service) AddDevice(payload *DevicePayload) (*Device, error) {
	err := s.validate(payload)
	if err != nil {
		return nil, err
	}
	return s.Dao.AddDevice(payload)
}
