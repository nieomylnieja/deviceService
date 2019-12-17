package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type HandlersEnv struct {
	controller HandlersController
}

func NewHandlersEnv(controller HandlersController) *HandlersEnv {
	return &HandlersEnv{
		controller: controller,
	}
}

func (he *HandlersEnv) AddDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var devPayload DevicePayload

	err := json.NewDecoder(r.Body).Decode(&devPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := he.controller.AddDevice(r.Context(), &devPayload)
	if caseSwitchError(w, err) {
		return
	}

	he.writeObject(w, device)

}

func (he *HandlersEnv) GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	device, err := he.controller.GetDevice(r.Context(), id)
	if device == nil && err == mongo.ErrNoDocuments {
		fmt.Println("device was not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if caseSwitchError(w, err) {
		return
	}

	he.writeObject(w, device)
}

func (he *HandlersEnv) GetPaginatedDevices(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value("limit").(int)
	page := r.Context().Value("page").(int)

	devices, err := he.controller.GetPaginatedDevices(r.Context(), limit, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	he.writeObject(w, devices)

}

func (he *HandlersEnv) StartTickerService(w http.ResponseWriter, r *http.Request) {
	err := he.controller.StartMeasurementService(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (he *HandlersEnv) writeObject(w http.ResponseWriter, object interface{}) {
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

func caseSwitchError(w http.ResponseWriter, err error) bool {
	if err != nil {
		switch err.(type) {
		case ErrValidation:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return true
	}
	return false
}
