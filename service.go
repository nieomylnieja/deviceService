package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type DevicePayload struct {
	Name     string  `json:"name" validate:"required,min=2,max=30"`
	Interval int     `json:"interval" validate:"gt=0,numeric"`
	Value    float64 `json:"value" validate:"numeric"`
}

type DeviceDao interface {
	AddDevice(device *DevicePayload) (int, error)
}

type Service struct {
	Dao       DeviceDao
	validator *validator.Validate
}

func (s *Service) run() {
	s.validator = validator.New()
}

func (s *Service) validate(payload *DevicePayload) error {
	validationErrors := s.validator.Struct(payload)
	if validationErrors != nil {
		for _, err := range validationErrors.(validator.ValidationErrors) {
			fmt.Println(err)
		}
		return NewErrValidation("input validation failed")
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

	d := Device{
		Id:       id,
		Name:     payload.Name,
		Value:    payload.Value,
		Interval: payload.Interval,
	}

	return &d, err
}
