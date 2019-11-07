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
	stopChan     chan bool
}

var devices map[int]Device
var devicesReadings map[int][]DeviceReading
var devicesChan chan DeviceInfo

func (d *Device) addDevice() {
	devices[d.ID] = *d
}

func (d *Device) updateDeviceName(name string) {
	d.Name = name
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
		case <-d.stopChan:
			ticker.Stop()
			fmt.Printf("...%s ID:%d stopped!\n", d.Name, d.ID)
			return
		}
	}
}

func (d *Device) stopDevice() {
	fmt.Printf("Stopping %s ID:%d...\n", d.Name, d.ID)
	d.stopChan <- true
	//delete(devices, d.ID)
}

func (d Device) createDevice(id int, name string, value string, interval float64) error {
	err := deviceAlreadyExistsCheck(id)
	if err != nil { return err }
	d = Device{id, name, value, interval, make(chan bool)}
	d.startDevice()
	return nil
}

func (d *Device) startDevice() () {
	d.addDevice()
	go d.deviceTicker()
}

func (d *Device) removeDevice() {
	delete(devices, d.ID)
	d.stopDevice()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%s ID:%d removed.\n", d.Name, d.ID)
}

func getDeviceByID(id int) *Device {
	dev := devices[id]
	return &dev
}

func deviceAlreadyExistsCheck(id int) error {
	var err error
	for _, v := range devices {
		if v.ID == id {
			err = errors.New("the device already exists")
			return err
		}
	}
	return nil
}

func (s *Service) tickerService() {
	var temp DeviceInfo
	var devRead DeviceReading
	for {
		select {
		case <-s.stopChan:
			return
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
	fmt.Fprint(w, "Here you will be shown the readings:\n")
	publishReadings(w)
}

type Service struct {
	stopChan chan bool
}

func (s *Service) run() {
	go s.tickerService()
}

func (s *Service) stop() {
	s.stopChan <- true
}

func main() {
	s := Service{make(chan bool)}
	devices = make(map[int]Device)
	devicesChan = make(chan DeviceInfo, 5)
	devicesReadings = make(map[int][]DeviceReading)

	s.run()

	var err error
	var dev Device
	for i := 0; i < 3; i++ {
		err = dev.createDevice(i, "Thermostat", "NULL", 1000)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	time.Sleep(2 * time.Second)

	for _, dev := range devices {
		dev.removeDevice()
	}

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}