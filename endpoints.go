package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/devices", devicesHandler).Methods("POST")
	router.HandleFunc("/", indexHandler).Methods("GET")

	return router
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

// TODO figure out how to actually return and read stuff properly
func devicesHandler(w http.ResponseWriter, r *http.Request) {
	resp := r.
		w.Write()
}
