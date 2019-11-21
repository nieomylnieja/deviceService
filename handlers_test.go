package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	defer resp.Body.Close()

	expected := &Device{Name: "test name"}

	var result *Device
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_GetAllDevicesHandler_GivenWrongInput_HandlerReturns400(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	urls := []string{"/devices?limit=a", "/devices?page=a",
		"/devices?limit=1&page=a", "/devices?limit=-1",
		"/devices?page=-1"}
	for _, url := range urls {
		resp, err := http.Get(mockServer.URL + url)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}

func Test_GetAllDevicesHandler_NoParams_HandlerDefaultsLimitAndPage(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_GetAllDevicesHandler_GivenDaoError_HandlerReturns500(t *testing.T) {
	r := newRouter(NewService(&mockDao{returnErr: ErrDao("")}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GetAllDevicesHandler_GivenLimitZero_HandlerReturnsAllDevices(t *testing.T) {
	m := &mockDao{data: make(map[int]Device)}
	m.data[0] = Device{Name: "test name"}
	r := newRouter(NewService(m))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices?limit=0")
	assert.NoError(t, err)
	defer resp.Body.Close()

	var expected []Device
	expected = append(expected, Device{Name: "test name"})

	var result []Device
	err = json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_GetAllDevicesHandler_GivenPageThatHasNoDevicesToShow_HandlerReturnsEmptyJsonArray(t *testing.T) {
	r := newRouter(NewService(&mockDao{}))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices?limit=1&page=1")
	assert.NoError(t, err)
	defer resp.Body.Close()

	expected := []int{}

	var result []int
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_GetAllDevicesHandler_GivenLimitGreaterThanSizeOfList_HandlerReturnsSuccess(t *testing.T) {
	m := &mockDao{data: make(map[int]Device)}
	m.data[0] = Device{Name: "test name"}
	r := newRouter(NewService(m))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_GetAllDevicesHandler_GivenLimitAndPage_HandlerReturnsCorrectList(t *testing.T) {
	m := &mockDao{data: make(map[int]Device)}
	for _, i := range []int{5, 3, 0, 25, 8} {
		m.data[i] = Device{Id: i}
	}
	r := newRouter(NewService(m))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices?limit=2&page=2")
	assert.NoError(t, err)
	defer resp.Body.Close()

	var expected []Device
	expected = append(expected, Device{Id: 25})

	var result []Device
	err = json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
