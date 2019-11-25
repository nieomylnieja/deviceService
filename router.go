package main

import (
	"github.com/gorilla/mux"
)

func newRouter(s *Service) *mux.Router {
	router := mux.NewRouter()

	devicesHandlerEnv := DeviceHandlers{s}
	router.HandleFunc("/devices", devicesHandlerEnv.AddDeviceHandler).Methods("POST")
	router.HandleFunc("/devices", pageAndLimitWrapper(devicesHandlerEnv.GetPaginatedDevices)).Methods("GET")
	router.HandleFunc("/devices/{id}", devicesHandlerEnv.GetDeviceHandler).Methods("GET")

	return router
}
