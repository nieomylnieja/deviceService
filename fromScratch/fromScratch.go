package main

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"time"
)

type DeviceInfo struct {
	ID    int
	Value string
	When  time.Time
}

type DeviceReading struct {
	Value string
	When  time.Time
}

type Device struct {
	ID       int
	Name     string
	Value    string
	Interval float64
	stop     chan bool
}

var devices map[int]Device
var devicesReadings map[int][]DeviceReading
var devicesChan chan DeviceInfo

func addDevice(d *Device) {
	devices[d.ID] = *d
}

func (d *Device) updateDeviceInterval(interval float64) {
	d.Interval = interval
}

func (d *Device) updateDeviceValue(value string) {
	d.Value = value
}

func (d *Device) deviceTicker() {
	ticker := time.NewTicker(time.Duration(d.Interval) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(d.Interval) * time.Millisecond)
			d.updateDeviceValue(valueService(d.ID))
			devicesChan <- DeviceInfo{d.ID, d.Value, time.Now()}
		case <-d.stop:
			ticker.Stop()
			fmt.Printf("...%s ID:%d stopped!\n", d.Name, d.ID)
			return
		}
	}
}

func stopDevice(d *Device) {
	fmt.Printf("Stopping %s ID:%d...\n", d.Name, d.ID)
	d.stop <- true
	//delete(devices, d.ID)
}

func createDevice(id int, name string, value string, interval float64) error {
	var err error
	for _, v := range devices {
		if v.ID == id {
			err = errors.New("the device already exists")
			return err
		}
	}
	d := Device{id, name, value, interval, make(chan bool)}
	d.startDevice()
	return nil
}

func (d *Device) startDevice() () {
	addDevice(d)
	go d.deviceTicker()
}

func removeDevice(d *Device) {
	delete(devices, d.ID)
	stopDevice(d)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%s ID:%d removed.\n", d.Name, d.ID)
}

func tickerService() {
	var temp DeviceInfo
	var devRead DeviceReading
	for {
		select {
		case temp = <-devicesChan:
			devRead = DeviceReading{temp.Value, temp.When}
			devicesReadings[temp.ID] = append(devicesReadings[temp.ID], devRead)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func publishReadings(w http.ResponseWriter) {
	for device, readings := range devicesReadings {
		fmt.Fprintf(w, "Device ID:%d\n", device)
		for _, r := range readings {
			fmt.Fprintf(w, "Nanoseconds: %d -- with value %s\n", r.When.Nanosecond(), r.Value)
		}
	}
}

// reads through the temperatures provided by lm-sensors
// giving every device a different reading
// on my laptop it provides 3 readings thus 3 devices
func valueService(n int) string {
	out, _ := exec.Command("sensors").Output()
	regexResult := regexp.MustCompile(".{4}.C\\s").FindAll(out, 6)
	result := string(regexResult[n])
	return result
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf("Here you will be shown the readings:")
	publishReadings(w)
}

func main() {
	devices = make(map[int]Device)
	devicesChan = make(chan DeviceInfo, 5)
	devicesReadings = make(map[int][]DeviceReading)

	go tickerService()

	var err error
	for i := 0; i < 3; i++ {
		err = createDevice(i, "Thermostat", "NULL", 1000)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	time.Sleep(10 * time.Second)

	for _, dev := range devices {
		removeDevice(&dev)
	}

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}
