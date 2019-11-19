package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*func Test_GivenTickerServiceIsStarted_WhenNewMeasurementComes_ThenReadingIsPassedToSaveChannel(t *testing.T) {
	t.Skip()
	s := Service{Dao: &Dao{Readings: make(map[int][]DeviceReading), Devices: make(map[int]Device)}}
	s.run()

	input := &RawInput{
		Id:       "0",
		Name:     "TestDevice",
		Interval: "1000",
	}
	devPayload, _ := s.CreateDevicePayload(input)
	m := func(n int) float64 { return 10.24 }
	dev, _ := s.Dao.AddDevice(devPayload)
	s.StartDevice(dev, m)

	time.Sleep(1 * time.Second)
	s.stop()

	expected := 10.24
	result := fmt.Sprintf("%v", s.dao.Readings[0][0].Value)

	assert.Equal(t, expected, result)
}*/

type mockDao struct {
	returnValue int
	returnErr   error
	calledTimes int
}

func (m *mockDao) AddDevice(device *DevicePayload) (int, error) {
	m.calledTimes++
	m.returnValue++
	return m.returnValue, m.returnErr
}

func Test_CorrectDevice_ServiceSavesNewDevice(t *testing.T) {
	device := &DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	}
	dao := &mockDao{returnValue: 1}

	out := Service{Dao: dao}

	result, err := out.AddDevice(device)

	assert.NoError(t, err)
	assert.Equal(t, 1, dao.calledTimes)
	assert.NotNil(t, result)
}

func Test_CorrectDeviceAndDaoFails_ServiceFails(t *testing.T) {
	dao := &mockDao{returnErr: fmt.Errorf("test error")}
	out := Service{Dao: dao}

	_, err := out.AddDevice(&DevicePayload{
		Value:    10.23,
		Name:     "Thermostat",
		Interval: 1000,
	})

	assert.Error(t, err)
}

func Test_GivenIntervalValueBelowZeroOrEqualToZero_ServiceFails(t *testing.T) {
	out := Service{Dao: &mockDao{}}

	_, err1 := out.AddDevice(&DevicePayload{Interval: -1})
	_, err2 := out.AddDevice(&DevicePayload{Interval: 0})
	assert.Error(t, err1)
	assert.Error(t, err2)
}

func Test_CorrectPayload_DaoFillsOutId_ServiceDefaultsInterval(t *testing.T) {
	dao := &mockDao{}
	out := Service{Dao: dao}

	_, err := out.AddDevice(&DevicePayload{Name: "aaa"})

	expected := &Device{
		Interval: 1000,
		Name:     "aaa",
		Id:       1,
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, out.tempDevice)
}

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

	dp := DevicePayload{"test name", 2, 0}
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
