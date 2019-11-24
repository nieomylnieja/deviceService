package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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

	id, err := convertToPositiveInteger(input)
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

func (dh *DeviceHandlers) GetManyDevicesHandler(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value("limit").(int)
	page := r.Context().Value("page").(int)

	devices, err := dh.service.GetManyDevices(limit, page)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dh.writeObject(w, devices)

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

func pageAndLimitWrapper(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := readIntFromQueryParameter(r.URL, "limit", 100)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		page, err := readIntFromQueryParameter(r.URL, "page", 100)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for key, val := range map[string]int{"limit": limit, "page": page} {
			r = r.WithContext(context.WithValue(r.Context(), key, val))
		}

		h.ServeHTTP(w, r)
	}
}
