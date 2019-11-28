package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

type HandlersEnvironment struct {
	controller *Controller
	startOnce  sync.Once
}

func NewHandlersEnvironment(controller *Controller) HandlersEnvironment {
	return HandlersEnvironment{
		controller: controller,
		startOnce:  sync.Once{},
	}
}

func (he *HandlersEnvironment) AddDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var devPayload DevicePayload

	err := json.NewDecoder(r.Body).Decode(&devPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := he.controller.AddDevice(&devPayload)
	if err != nil {
		switch err.(type) {
		case ErrValidation:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	he.writeObject(w, device)

}

func (he *HandlersEnvironment) GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	input := mux.Vars(r)["id"]

	id, err := convertToPositiveInteger(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := he.controller.GetDevice(id)
	if device == nil && err == nil {
		fmt.Println("device was not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	he.writeObject(w, device)

}

func (he *HandlersEnvironment) GetPaginatedDevices(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value("limit").(int)
	page := r.Context().Value("page").(int)

	devices, err := he.controller.GetPaginatedDevices(limit, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	he.writeObject(w, devices)

}

func (he *HandlersEnvironment) StartTickerService(w http.ResponseWriter, r *http.Request) {
	he.startOnce.Do(func() {
		err := he.controller.StartTickerService()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (he *HandlersEnvironment) writeObject(w http.ResponseWriter, object interface{}) {
	respBody, err := json.Marshal(object)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func pageAndLimitWrapper(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := readIntFromQueryParameter(r.URL, "limit", 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		page, err := readIntFromQueryParameter(r.URL, "page", 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "limit", limit)
		ctx = context.WithValue(ctx, "page", page)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
