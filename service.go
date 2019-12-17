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
	Dao       DevicesDao
	validator *validator.Validate
}

func NewService(dao DevicesDao) *Service {
	return &Service{
		Dao:       dao,
		validator: validator.New(),
	}
}

func (s *Service) AddDevice(ctx context.Context, payload *DevicePayload) (*Device, error) {
	if payload.Interval == 0 {
		payload.Interval = 1000
	}
	if err := s.validateDevicePayload(payload); err != nil {
		return nil, err
	}
	id, err := s.Dao.AddDevice(ctx, payload)
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

func (s *Service) GetDevice(ctx context.Context, id string) (*Device, error) {
	objectID, err := stringIDToObjectID(id)
	if err != nil {
		return nil, ErrValidation("")
	}
	return s.Dao.GetDevice(ctx, objectID)
}

func (s *Service) GetPaginatedDevices(ctx context.Context, limit, page int) ([]Device, error) {
	return s.Dao.GetPaginatedDevices(ctx, limit, page)
}

func (s *Service) GetAllDevices(ctx context.Context) ([]Device, error) {
	return s.Dao.GetAllDevices(ctx)
}

func (s *Service) validateDevicePayload(payload *DevicePayload) error {
	validationErrors := s.validator.Struct(payload)
	if validationErrors != nil {
		for _, err := range validationErrors.(validator.ValidationErrors) {
			fmt.Println(err)
		}
		return ErrValidation("")
	}
	return nil
}
