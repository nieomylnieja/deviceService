package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WriteObject_GivenAnObject_FuncWritesMarshalledObject(t *testing.T) {
	dh := DeviceHandlers{}
	resp := httptest.NewRecorder()

	dh.writeObject(resp, Device{Id: 1})

	var actual Device
	err := json.NewDecoder(resp.Body).Decode(&actual)

	assert.NoError(t, err)
	assert.Equal(t, Device{Id: 1}, actual)
}

func Test_ConvertToPositiveInteger_GivenWrongInput_FuncReturnsError(t *testing.T) {
	dh := DeviceHandlers{}

	tests := map[string]string{
		"char":            "a",
		"negative number": "-2",
		"float":           "0.0",
		"interface":       "{}",
	}

	for name, tc := range tests {
		_, err := dh.convertToPositiveInteger(tc)
		assert.Error(t, err, name)
	}
}

func Test_ConvertToPositiveInteger_GivenCorrectInput_FuncReturnsPositiveInt(t *testing.T) {
	dh := DeviceHandlers{}

	tests := map[string]string{
		"zero":              "0",
		"non zero positive": "14",
	}

	for name, tc := range tests {
		actual, err := dh.convertToPositiveInteger(tc)

		assert.NoError(t, err, name)
		assert.IsType(t, 1, actual, name)
		assert.GreaterOrEqual(t, actual, 0, name)
	}
}

func Test_AddDeviceHandler_GivenInvalidDevicePayload_HandlerReturns400(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test", "interval": -1}`))
	resp, _ := http.Post(mockServer.URL+"/devices", "application/json", requestBody)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_AddDeviceHandler_GivenDevicePayload_HandlerReturnsDeviceObjectAndPerformsAddDevice(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	dp := DevicePayload{Name: "test name", Interval: 2}
	requestBody, err := json.Marshal(dp)
	assert.NoError(t, err)

	resp, err := http.Post(mockServer.URL+"/devices", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

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

func Test_GetDeviceHandler_GivenNonNumericId_HandlerReturnsError400(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices/test")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_GetDeviceHandler_GivenNonExistingId_HandlerReturnsError404(t *testing.T) {
	r := newRouter(NewService(&mockDao{device: nil}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices/123")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_GetDeviceHandler_GivenErrorInDao_HandlerReturnsError500(t *testing.T) {
	r := newRouter(NewService(&mockDao{returnErr: ErrDao("")}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices/1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GetDeviceHandler_GivenCorrectId_HandlerReturnsDeviceObject(t *testing.T) {
	r := newRouter(NewService(&mockDao{device: &Device{Name: "test name"}}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices/1")
	assert.NoError(t, err)

	expected := &Device{Name: "test name"}

	var result *Device
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_GetManyDevicesHandler_GivenWrongInput_HandlerReturns400(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	tests := map[string]string{
		"limit is not int":               "/devices?limit=a",
		"page is not int":                "/devices?page=a",
		"page is not int, correct limit": "/devices?limit=1&page=a",
		"limit is below zero":            "/devices?limit=-1",
		"page is below zero":             "/devices?page=-1",
	}
	for name, tc := range tests {
		resp, err := http.Get(mockServer.URL + tc)

		assert.NoError(t, err, name)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, name)
	}
}

func Test_GetManyDevicesHandler_NoParams_HandlerDefaultsLimitAndPage(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_GetManyDevicesHandler_GivenDaoError_HandlerReturns500(t *testing.T) {
	r := newRouter(NewService(&mockDao{returnErr: ErrDao("")}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GetManyDevicesHandler_GivenLimitZero_HandlerReturnsAllDevices(t *testing.T) {
	r := newRouter(NewService(&mockDao{data: []Device{{Name: "test name"}}}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices?limit=0")
	assert.NoError(t, err)

	expected := []Device{{Name: "test name"}}

	var result []Device
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_GetManyDevicesHandler_GivenPageThatHasNoDevicesToShow_HandlerReturnsEmptyJsonArray(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices?page=1")
	assert.NoError(t, err)

	var result []int
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, []int(nil), result)
}
