package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GivenNonExistingRoute_RouterReturns404(t *testing.T) {
	dao := &mockDao{}
	out := NewService(dao)
	r := newRouter(out)
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/dcisve")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_GivenInvalidMethod_RouterReturns405(t *testing.T) {
	dao := &mockDao{}
	out := NewService(dao)
	r := newRouter(out)
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func Test_GivenDaoError_RouterReturns500(t *testing.T) {
	dao := &mockDao{returnErr: ErrDao("")}
	out := NewService(dao)
	r := newRouter(out)
	mockServer := httptest.NewServer(r)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test"}`))
	resp, _ := http.Post(mockServer.URL+"/devices", "application/json", requestBody)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GivenInvalidDevicePayload_HandlerReturns400(t *testing.T) {
	dao := &mockDao{}
	out := NewService(dao)
	r := newRouter(out)
	mockServer := httptest.NewServer(r)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test", "interval": -1}`))
	resp, _ := http.Post(mockServer.URL+"/devices", "application/json", requestBody)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_GivenDevicePayload_HandlerReturnsDeviceObjectAndPerformsAddDevice(t *testing.T) {
	dao := &mockDao{}
	out := NewService(dao)
	r := newRouter(out)
	mockServer := httptest.NewServer(r)

	dp := DevicePayload{Name: "test name", Interval: 2}
	requestBody, err := json.Marshal(dp)
	assert.NoError(t, err)

	resp, err := http.Post(mockServer.URL+"/devices", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	expected := Device{
		Id:       1,
		Name:     "test name",
		Value:    0,
		Interval: 2,
	}

	var result Device

	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
