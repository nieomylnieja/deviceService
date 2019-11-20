package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeviceHandlers struct {
	service *Service
}

func (dh *DeviceHandlers) addDeviceHandler(w http.ResponseWriter, r *http.Request) {
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

	return
}

func (dh *DeviceHandlers) getDeviceHandler(w http.ResponseWriter, r *http.Request) {
	return
}
