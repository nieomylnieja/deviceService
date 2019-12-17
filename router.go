package main

import (
	"context"
	"github.com/gorilla/mux"
)

type HandlersController interface {
	StartMeasurementService(ctx context.Context) error
	GetDevice(ctx context.Context, id string) (*Device, error)
	AddDevice(ctx context.Context, devPayload *DevicePayload) (*Device, error)
	GetPaginatedDevices(ctx context.Context, limit, page int) ([]Device, error)
}

func NewRouter(hCtrl HandlersController) *mux.Router {
	router := mux.NewRouter()

	hEnv := NewHandlersEnv(hCtrl)
	router.HandleFunc("/start", hEnv.StartTickerService).Methods("POST")
	router.HandleFunc("/devices", hEnv.AddDeviceHandler).Methods("POST")
	router.HandleFunc("/devices", pageAndLimitWrapper(hEnv.GetPaginatedDevices)).Methods("GET")
	router.HandleFunc("/devices/{id}", hEnv.GetDeviceHandler).Methods("GET")

	return router
}
