package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WriteObject_GivenAnObject_FuncWritesMarshalledObject(t *testing.T) {
	dh := HandlersEnvironment{}
	resp := httptest.NewRecorder()
	id := primitive.NewObjectID()

	dh.writeObject(resp, Device{Id: id})

	var actual Device
	err := json.NewDecoder(resp.Body).Decode(&actual)

	assert.NoError(t, err)
	assert.Equal(t, Device{Id: id}, actual)
}

func Test_AddDeviceHandler_GivenInvalidDevicePayload_HandlerReturns400(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{})})
	mockServer := httptest.NewServer(r)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test", "interval": -1}`))
	resp, _ := http.Post(mockServer.URL+"/devices", "application/json", requestBody)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_AddDeviceHandler_GivenDevicePayload_HandlerReturnsDeviceObjectAndPerformsAddDevice(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{})})
	mockServer := httptest.NewServer(r)

	dp := DevicePayload{Name: "test name", Interval: 2}
	requestBody, err := json.Marshal(dp)
	assert.NoError(t, err)

	resp, err := http.Post(mockServer.URL+"/devices", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	expected := Device{
		Name:     "test name",
		Value:    0,
		Interval: 2,
	}

	var result Device
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_GetDeviceHandler_GivenNonExistingId_HandlerReturnsError404(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{device: nil})})
	mockServer := httptest.NewServer(r)
	id := primitive.NewObjectID().Hex()

	resp, err := http.Get(mockServer.URL + "/devices/" + id)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_GetDeviceHandler_GivenErrorInDao_HandlerReturnsError500(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{returnErr: ErrDao("")})})
	mockServer := httptest.NewServer(r)
	id := primitive.NewObjectID().Hex()

	resp, err := http.Get(mockServer.URL + "/devices/" + id)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GetDeviceHandler_GivenCorrectId_HandlerReturnsDeviceObject(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{device: &Device{Name: "test name"}})})
	mockServer := httptest.NewServer(r)
	id := primitive.NewObjectID().Hex()

	resp, err := http.Get(mockServer.URL + "/devices/" + id)
	assert.NoError(t, err)

	expected := &Device{Name: "test name"}

	var result *Device
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_PageAndLimitWrapper_GivenWrongInput_HandlerReturns400(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{})})
	mockServer := httptest.NewServer(r)

	tests := map[string]string{
		"limit is not int":               "/devices?limit=a",
		"page is not int":                "/devices?page=a",
		"page is not int, correct limit": "/devices?limit=1&page=a",
		"limit is below zero":            "/devices?limit=-1",
		"page is below zero":             "/devices?page=-1",
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Get(mockServer.URL + tc)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func Test_PageAndLimitWrapper_NoParams_HandlerDefaultsLimitAndPage(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{})})
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_GetPaginatedDevicesHandler_GivenDaoError_HandlerReturns500(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{returnErr: ErrDao("")})})
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GetPaginatedDevicesHandler_GivenPageThatHasNoDevicesToShow_HandlerReturnsEmptyJsonArray(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{})})
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/devices?page=1")
	assert.NoError(t, err)

	var result []int
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, []int(nil), result)
}

func Test_StartTickerServiceHandler_GivenDaoError_HandlerReturns500(t *testing.T) {
	r := newRouter(&Controller{mainService: NewService(&mockDao{returnErr: ErrDao("")})})
	mockServer := httptest.NewServer(r)

	resp, err := http.Post(mockServer.URL+"/start", "", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_CaseSwitchError_GivenDifferentErrors_FuncWritesProperStatusCode(t *testing.T) {
	tests := map[string]struct {
		err      error
		expected int
	}{
		"validator err": {err: ErrValidation(""), expected: http.StatusBadRequest},
		"dao err":       {err: ErrDao(""), expected: http.StatusInternalServerError},
	}

	var err bool
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			w := &mockResponseWriter{}
			err = caseSwitchError(w, tc.err)

			assert.Equal(t, err, true)
			assert.Equal(t, tc.expected, w.calledWithStatusCode)
		})
	}
}

type mockResponseWriter struct {
	calledWithStatusCode int
}

func (m *mockResponseWriter) Header() http.Header {
	return map[string][]string{}
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.calledWithStatusCode = statusCode
}
