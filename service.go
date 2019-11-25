package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type DevicePayload struct {
	Name     string  `json:"name" validate:"required,min=2,max=30"`
	Interval int     `json:"interval,string" validate:"gt=0,numeric"`
	Value    float64 `json:"value,string" validate:"numeric"`
}

type DeviceDao interface {
	AddDevice(device *DevicePayload) (int, error)
	GetDevice(id int) (*Device, error)
	GetPaginatedDevices(limit int, page int) ([]Device, error)
	GetAllDevices() ([]Device, error)
}

type Service struct {
	Dao       DeviceDao
	validator *validator.Validate
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
	return s.Dao.GetDevice(id)
}

func (s *Service) GetPaginatedDevices(limit int, page int) ([]Device, error) {
	return s.Dao.GetPaginatedDevices(limit, page)
}

func (s *Service) StartTickerService() error {
	devices, err := s.Dao.GetAllDevices()
	if err != nil {
		return err
	}
	mch := make(chan Measurement)
	quit := make(chan bool)
	defer close(quit)
	defer close(mch)
	go s.MeasurementsWriter(mch)
	for i := range devices {
		go devices[i].deviceTicker(mch, quit)
	}
	select {}
}

func (s *Service) MeasurementsWriter(mch chan Measurement) {
	for m := range mch {
		fmt.Printf("ID:%d -- %f\n", m.Id, m.Value)
	}
}
