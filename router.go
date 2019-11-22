package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func newRouter(s *Service) *mux.Router {
	router := mux.NewRouter()

	devicesHandlerEnv := DeviceHandlers{s}
	router.HandleFunc("/devices", devicesHandlerEnv.AddDeviceHandler).Methods("POST")
	router.HandleFunc("/devices", pageAndLimitWrapper(devicesHandlerEnv.GetManyDevicesHandler)).Methods("GET")
	router.HandleFunc("/devices/{id}", devicesHandlerEnv.GetDeviceHandler).Methods("GET")

	return router
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
