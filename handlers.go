package main

import (
	"context"
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
		var err error
		var limit, page int

		limitStr := r.URL.Query().Get("limit")
		if limitStr == "" {
			limit = 100
		} else {
			limit, err = convertToPositiveInteger(limitStr)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			page = 0
		} else {
			page, err = convertToPositiveInteger(pageStr)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		for key, val := range map[string]int{"limit": limit, "page": page} {
			r = r.WithContext(context.WithValue(r.Context(), key, val))
		}

		h.ServeHTTP(w, r)
	}
}

func convertToPositiveInteger(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if id < 0 {
		return 0, errors.New("input is a negative number")
	}
	return id, nil
}
