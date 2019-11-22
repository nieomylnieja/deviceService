package main

import (
	"encoding/json"
	"errors"
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

}

func (dh *DeviceHandlers) GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	input := mux.Vars(r)["id"]

	id, err := dh.convertToPositiveInteger(input)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
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

}

func (dh *DeviceHandlers) GetAllDevicesHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var limit, page int
	keys := r.URL.Query()

	limitStr := keys.Get("limit")
	if limitStr == "" {
		limit = 100
	} else {
		limit, err = dh.convertToPositiveInteger(limitStr)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	pageStr := keys.Get("page")
	if pageStr == "" {
		page = 0
	} else {
		page, err = dh.convertToPositiveInteger(pageStr)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	devices, err := dh.service.GetSortedDevicesList()
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if limit == 0 {
		dh.writeObject(w, devices)
		return
	}

	if limit*page > len(*devices) {
		dh.writeObject(w, nil)
		return
	}

	if page*limit+limit >= len(*devices) {
		dh.writeObject(w, (*devices)[page*limit:len(*devices)])
		return
	}

	dh.writeObject(w, (*devices)[page*limit:page*limit+limit])

}

func (dh *DeviceHandlers) writeObject(w http.ResponseWriter, object interface{}) {
	respBody, err := json.Marshal(object)
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

func (dh *DeviceHandlers) convertToPositiveInteger(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if id < 0 {
		return 0, errors.New("input is a negative number")
	}
	return id, nil
}
