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

	id, ok := dh.stringIsPositiveNumberReturnInt(w, input)
	if !ok {
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
	var ok bool
	var limit, page int
	keys := r.URL.Query()

	limitStr := keys.Get("limit")
	if limitStr == "" {
		limit = 100
	} else {
		limit, ok = dh.stringIsPositiveNumberReturnInt(w, limitStr)
		if !ok {
			return
		}
	}
	pageStr := keys.Get("page")
	if pageStr == "" {
		page = 0
	} else {
		page, ok = dh.stringIsPositiveNumberReturnInt(w, pageStr)
		if !ok {
			return
		}
	}

	fmt.Println(limit, page)
	/*devices, err := dh.service.GetAllDevices()
	if err != nil {
		return
	}

	if limit == 0 {
		respBody, _ := json.Marshal(*devices)
		_, _ = w.Write(respBody)
	}

	if len(*devices) / limit > page {
		emptyArr := make(map[int]Device, 0)
		respBody, _ := json.Marshal(emptyArr)
		_, _ = w.Write(respBody)
	}
	*/
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

func (dh *DeviceHandlers) stringIsPositiveNumberReturnInt(w http.ResponseWriter, input string) (int, bool) {
	id, err := strconv.Atoi(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return 0, false
	}
	if id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(errors.New("input is a negative number"))
	}
	return id, true
}
