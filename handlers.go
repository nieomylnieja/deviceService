package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeviceHandlers struct {
	service *Service
}

func (h *DeviceHandlers) devicesHandler(w http.ResponseWriter, r *http.Request) {
	var devPayload DevicePayload

	err := json.NewDecoder(r.Body).Decode(&devPayload)
	if err != nil {
		fmt.Printf("handlerError: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	device, err := h.service.AddDevice(&devPayload)
	if err != nil {
		switch err.(type) {
		case *ErrDao:
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		case *ErrValidation:
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
		default:
			fmt.Println("unhandled error!")
			w.WriteHeader(http.StatusNotImplemented)
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
