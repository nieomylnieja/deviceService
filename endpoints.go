package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
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
	var devPayload DevicePayload

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &devPayload)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d := Device{
		Id:       1,
		Name:     devPayload.Name,
		Value:    devPayload.Value,
		Interval: devPayload.Interval,
		stopChan: nil,
	}

	respBody, err := json.Marshal(d)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
}
