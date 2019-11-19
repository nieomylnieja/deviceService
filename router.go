package main

import (
	"github.com/gorilla/mux"
)

type RouterEnv struct {
	service *Service
}

func (r *RouterEnv) newRouter() *mux.Router {
	router := mux.NewRouter()

	indexHandlerEnv := HandlersEnv{r.service}
	indexHandler := indexHandlerEnv.indexHandler
	devicesHandlerEnv := HandlersEnv{r.service}
	devicesHandler := devicesHandlerEnv.devicesHandler
	router.HandleFunc("/devices", devicesHandler).Methods("POST")
	router.HandleFunc("/", indexHandler).Methods("GET")

	return router
}
