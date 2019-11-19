package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GivenCorrectMethodAndRoute_RouterSuccessAndBodyMatches(t *testing.T) {
	renv := RouterEnv{}
	r := renv.newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/")

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	expected := "hello"

	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))
}

func Test_GivenNonExistingRoute_RouterReturns404(t *testing.T) {
	renv := RouterEnv{}
	r := renv.newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/dcisve")

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func Test_GivenInvalidMethod_RouterReturns405(t *testing.T) {
	renv := RouterEnv{}
	r := renv.newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Post(mockServer.URL+"/", "", nil)

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)
}

func Test_GivenDevicePayload_HandlerReturnsDeviceObjectAndPerformsAddDevice(t *testing.T) {
	dao := &mockDao{}
	out := Service{Dao: dao}

	renv := RouterEnv{&out}
	r := renv.newRouter()
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
		stopChan: nil,
	}

	var result Device

	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_GivenInvalidDevicePayload_HandlerReturns400(t *testing.T) {
	dao := &mockDao{}
	out := Service{Dao: dao}

	renv := RouterEnv{&out}
	r := renv.newRouter()
	mockServer := httptest.NewServer(r)

	dp := DevicePayload{Name: "test name", Interval: -1}
	requestBody, err := json.Marshal(dp)
	assert.NoError(t, err)

	resp, _ := http.Post(mockServer.URL+"/devices", "application/json", bytes.NewBuffer(requestBody))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
