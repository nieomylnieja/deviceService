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

type ParsedInput struct {
	Id       int
	Name     string
	Interval int
}

type Service struct {
	devices         map[int]Device
	DevicesReadings *DataAccessObject
	DevicesSaveChan chan DeviceInfo
	stopChan        chan bool
}

func (s *Service) init(dao *DataAccessObject) {
	s.DevicesReadings = dao
	s.DevicesSaveChan = make(chan DeviceInfo)
	s.stopChan = make(chan bool)
	s.devices = make(map[int]Device)
}

func (s *Service) run() {
	go s.tickerService()
}

func (s *Service) stop() {
	for _, dev := range s.devices {
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
			s.DevicesReadings.data[deviceInfo.Id] =
				append(s.DevicesReadings.data[deviceInfo.Id], deviceReading)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (s *Service) parseDeviceInitInput(input *RawInput) (*ParsedInput, error) {
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
	parsedInput := &ParsedInput{id, input.Name, interval}
	return parsedInput, nil
}

func (s *Service) createDevice(input *RawInput) (*Device, error) {
	parsedInput, err := s.parseDeviceInitInput(input)
	if err != nil {
		return nil, err
	}
	d := &Device{parsedInput.Id, parsedInput.Name, "",
		parsedInput.Interval, make(chan bool)}
	return d, nil
}

func (s *Service) updateDeviceName(id int, name string) error {
	dev, err := s.getDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Name = name
	return nil
}

func (s *Service) updateDeviceInterval(id int, interval int) error {
	dev, err := s.getDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Interval = interval
	return nil
}

func (s *Service) updateDeviceValue(d *Device, value string) {
	d.Value = value
}

func (s *Service) startDevice(d *Device, getMeasurement measurement) error {
	err := s.addDevice(d)
	if err != nil {
		return err
	}
	go d.deviceTicker(s, getMeasurement)
	return nil
}

func (s *Service) addDevice(d *Device) error {
	if s.deviceAlreadyExists(d.Id) {
		err := errors.New("The device with given ID already exists!")
		return err
	}
	s.devices[d.Id] = *d
	return nil
}

func (s *Service) stopDevice(d *Device) {
	fmt.Printf("Stopping %s ID:%d...\n", d.Name, d.Id)
	d.stopChan <- true
	s.removeDevice(d)
}

func (s *Service) removeDevice(d *Device) {
	delete(s.devices, d.Id)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%s ID:%d removed.\n", d.Name, d.Id)
}

func (s *Service) getDevicesList() {
	for _, dev := range s.devices {
		fmt.Printf("%s -- ID:%d and interval=%d mls\n",
			dev.Name, dev.Id, dev.Interval)
	}
}

func (s *Service) getDeviceByID(id int) (*Device, error) {
	var err error
	if s.deviceAlreadyExists(id) {
		dev := s.devices[id]
		return &dev, nil
	}
	err = errors.New("The device with given ID doesn't exist!")
	return nil, err
}

func (s *Service) deviceAlreadyExists(id int) bool {
	for _, dev := range s.devices {
		if dev.Id == id {
			return true
		}
	}
	return false
}

func (s *Service) getReadings() []string {
	var fwdReadings []string
	for device, readings := range s.DevicesReadings.data {
		fwdReadings = append(fwdReadings, fmt.Sprintf("Device ID:%d\n", device))
		for _, r := range readings {
			fwdReadings = append(fwdReadings,
				fmt.Sprintf("Nanoseconds: %d -- with value %s\n", r.When.Nanosecond(), r.Value))
		}
	}
	return fwdReadings
}