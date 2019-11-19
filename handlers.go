package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HandlersEnv struct {
	service *Service
}

func (h *HandlersEnv) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func (h *HandlersEnv) devicesHandler(w http.ResponseWriter, r *http.Request) {
	var devPayload DevicePayload

	err := json.NewDecoder(r.Body).Decode(&devPayload)
	if err != nil {
		fmt.Println(fmt.Errorf("handlerError: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.service.AddDevice(&devPayload)
	if err != nil {
		fmt.Println(fmt.Errorf("handlerError: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d := h.service.tempDevice

	respBody, err := json.Marshal(d)
	if err != nil {
		fmt.Println(fmt.Errorf("handlerError: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(respBody)
	return
}
