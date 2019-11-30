package main

import (
	"github.com/gorilla/mux"
)

func newRouter(c *Controller) *mux.Router {
	router := mux.NewRouter()

	handlersEnvironment := NewHandlersEnvironment(c)
	router.HandleFunc("/start", handlersEnvironment.StartTickerService).Methods("POST")
	router.HandleFunc("/devices", handlersEnvironment.AddDeviceHandler).Methods("POST")
	router.HandleFunc("/devices", pageAndLimitWrapper(handlersEnvironment.GetPaginatedDevices)).Methods("GET")
	router.HandleFunc("/devices/{id}", handlersEnvironment.GetDeviceHandler).Methods("GET")

	return router
}
