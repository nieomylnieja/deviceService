package main

import (
	"github.com/gorilla/mux"
)

func newRouter(s *Service) *mux.Router {
	router := mux.NewRouter()

	devicesHandlerEnv := DeviceHandlers{s}
	router.HandleFunc("/devices", devicesHandlerEnv.addDeviceHandler).Methods("POST")
	router.HandleFunc("/devices/{id}", devicesHandlerEnv.getDeviceHandler).Methods("GET")

	return router
}
