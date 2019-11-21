package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"sort"
)

type DevicePayload struct {
	Name     string  `json:"name" validate:"required,min=2,max=30"`
	Interval int     `json:"interval" validate:"gt=0,numeric"`
	Value    float64 `json:"value" validate:"numeric"`
}

type DeviceDao interface {
	AddDevice(device *DevicePayload) (int, error)
	GetDevice(id int) (*Device, error)
	GetAllDevices() (*map[int]Device, error)
}

type Service struct {
	Dao         DeviceDao
	validator   *validator.Validate
	indexHolder []int
}

func NewService(dao DeviceDao) *Service {
	s := Service{
		Dao:       dao,
		validator: validator.New(),
	}
	return &s
}

func (s *Service) validate(payload *DevicePayload) error {
	validationErrors := s.validator.Struct(payload)
	if validationErrors != nil {
		for _, err := range validationErrors.(validator.ValidationErrors) {
			fmt.Println(err)
		}
		return ErrValidation("")
	}
	return nil
}

func (s *Service) AddDevice(payload *DevicePayload) (*Device, error) {
	if payload.Interval == 0 {
		payload.Interval = 1000
	}
	err := s.validate(payload)
	if err != nil {
		return nil, err
	}
	id, err := s.Dao.AddDevice(payload)
	if err != nil {
		return nil, err
	}

	return &Device{
		Id:       id,
		Name:     payload.Name,
		Value:    payload.Value,
		Interval: payload.Interval,
	}, nil
}

func (s *Service) GetDevice(id int) (*Device, error) {
	device, err := s.Dao.GetDevice(id)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (s *Service) GetAllDevices() (*map[int]Device, error) {
	devices, err := s.Dao.GetAllDevices()
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *Service) GetSortedDevicesList() (*[]Device, error) {
	devices, err := s.GetAllDevices()
	if err != nil {
		return nil, err
	}

	var keys []int
	for k := range *devices {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	sortedDevicesList := []Device{} // I think I ought to initialize it so that JSON can actually return and empty list
	for _, k := range keys {
		sortedDevicesList = append(sortedDevicesList, (*devices)[k])
	}

	return &sortedDevicesList, nil
}
