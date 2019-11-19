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
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.service.AddDevice(&devPayload)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d := h.service.tempDevice

	respBody, err := json.Marshal(d)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
}
