package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type DeviceHandlers struct {
	service *Service
}

func (dh *DeviceHandlers) AddDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var devPayload DevicePayload

	err := json.NewDecoder(r.Body).Decode(&devPayload)
	if err != nil {
		fmt.Printf("handlerError: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	device, err := dh.service.AddDevice(&devPayload)
	if err != nil {
		switch err.(type) {
		case ErrValidation:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err.Error())
		default:
			fmt.Println("unhandled error!")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	dh.writeObject(w, device)

	return
}

func (dh *DeviceHandlers) GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	input := mux.Vars(r)["id"]
	id, err := strconv.Atoi(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	device, err := dh.service.GetDevice(id)
	if device == nil && err == nil {
		fmt.Println("device was not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dh.writeObject(w, device)

	return
}

func (dh *DeviceHandlers) GetAllDevicesHandler(w http.ResponseWriter, device *Device) {
	return
}

func (dh *DeviceHandlers) writeObject(w http.ResponseWriter, device *Device) {
	respBody, err := json.Marshal(device)
	if err != nil {
		fmt.Printf("handlerError: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		fmt.Printf("handlerError: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
