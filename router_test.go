package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GivenNonExistingRoute_RouterReturns404(t *testing.T) {
	r := newRouter(NewController(NewService(&mockDao{})))
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/dcisve")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_GivenInvalidMethod_RouterReturns405(t *testing.T) {
	r := newRouter(NewController(NewService(&mockDao{})))
	mockServer := httptest.NewServer(r)

	resp, err := http.Post(mockServer.URL+"/devices/2", "", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func Test_GivenDaoError_RouterReturns500(t *testing.T) {
	r := newRouter(NewController(NewService(&mockDao{returnErr: ErrDao("")})))
	mockServer := httptest.NewServer(r)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test"}`))
	resp, _ := http.Post(mockServer.URL+"/devices", "application/json", requestBody)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
