package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type DevicePayload struct {
	Name     string  `json:"name" validate:"required,min=2,max=30"`
	Interval int     `json:"interval,string" validate:"gt=0,numeric"`
	Value    float64 `json:"value,string" validate:"numeric"`
}

type Service struct {
	Dao       DeviceDao
	validator *validator.Validate
}

func NewService(dao DeviceDao) *Service {
	return &Service{
		Dao:       dao,
		validator: validator.New(),
	}
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

func (s *Service) AddDevice(payload *DevicePayload, ctx context.Context) (*Device, error) {
	if payload.Interval == 0 {
		payload.Interval = 1000
	}
	err := s.validate(payload)
	if err != nil {
		return nil, err
	}
	id, err := s.Dao.AddDevice(payload, ctx)
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

func (s *Service) GetDevice(id string, ctx context.Context) (*Device, error) {
	return s.Dao.GetDevice(id, ctx)
}

func (s *Service) GetPaginatedDevices(limit, page int, ctx context.Context) ([]Device, error) {
	return s.Dao.GetPaginatedDevices(limit, page, ctx)
}

func (s *Service) GetAllDevices(ctx context.Context) ([]Device, error) {
	return s.Dao.GetAllDevices(ctx)
}
