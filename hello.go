package main

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"time"
)

// DEVICE METHODS AND STRUCTS

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
	stopChan chan bool
}

func (d *Device) addDevice(s *Service) {
	s.devices[d.ID] = *d
}

func (d *Device) deviceTicker(s *Service) {
	ticker := time.NewTicker(time.Duration(d.Interval) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(d.Interval) * time.Millisecond)
			s.updateDeviceValue(d, valueService(d.ID))
			s.devicesChan <- DeviceInfo{d.ID, d.Value, time.Now()}
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
}

func (d *Device) startDevice(s *Service) () {
	go d.deviceTicker(s)
}

// SERVICE METHODS AND STRUCTS

type Service struct {
	devices         map[int]Device
	devicesChan     chan DeviceInfo
	stopChan        chan bool
}

func (s *Service) run() {
	devicesReadings = make(map[int][]DeviceReading)
	s.devicesChan = make(chan DeviceInfo)
	s.stopChan = make(chan bool)
	s.devices = make(map[int]Device)
	go s.tickerService()
}

func (s *Service) stop() {
	for _, dev := range s.devices {
		s.removeDevice(&dev)
	}
	s.stopChan <- true
}

func (s *Service) tickerService() {
	var temp DeviceInfo
	var devRead DeviceReading
	for {
		select {
		case <-s.stopChan:
			return
		case temp = <-s.devicesChan:
			devRead = DeviceReading{temp.Value, temp.When}
			devicesReadings[temp.ID] = append(devicesReadings[temp.ID], devRead)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (s *Service) createDevice(id int, name string, value string, interval float64) error {
	if s.deviceAlreadyExists(id) {
		err := errors.New("The device with given ID already exists!")
		return err
	}
	d := Device{id, name, value, interval, make(chan bool)}
	d.addDevice(s)
	d.startDevice(s)
	return nil
}

func (s *Service) updateDeviceName(id int, name string) error {
	dev, err := s.getDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Name = name
	return nil
}

func (s *Service) updateDeviceInterval(id int, interval float64) error {
	dev, err := s.getDeviceByID(id)
	if err != nil {
		return err
	}
	dev.Interval = interval
	return nil
}

func (s *Service) updateDeviceValue(d *Device, value string) {
	d.Value = value
}

func (s *Service) removeDevice(d *Device) {
	delete(s.devices, d.ID)
	d.stopDevice()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%s ID:%d removed.\n", d.Name, d.ID)
}

func (s *Service) getDevicesList() {
	for _, dev := range s.devices {
		fmt.Printf("%s -- ID:%d and interval=%f mls\n",
			dev.Name, dev.ID, dev.Interval)
	}
}

func (s *Service) getDeviceByID(id int) (*Device, error) {
	var err error
	if s.deviceAlreadyExists(id) {
		dev := s.devices[id]
		return &dev, nil
	}
	err = errors.New("The device with given ID doesn't exist!")
	return nil, err
}

func (s *Service) deviceAlreadyExists(id int) bool {
	for _, dev := range s.devices {
		if dev.ID == id {
			return true
		}
	}
	return false
}

// PERSISTENCE LAYER

var devicesReadings map[int][]DeviceReading

// reads through the temperatures provided by lm-sensors
// giving every device a different reading
// on my laptop it provides 3 readings thus 3 devices
func valueService(n int) string {
	out, _ := exec.Command("sensors").Output()
	regexResult := regexp.MustCompile(".{4}.C\\s").FindAll(out, 6)
	result := string(regexResult[n])
	return result
}

func pushReadings() []string {
	var fwdReadings []string
	for device, readings := range devicesReadings {
		fwdReadings = append(fwdReadings, fmt.Sprintf("Device ID:%d\n", device))
		for _, r := range readings {
			fwdReadings = append(fwdReadings,
				fmt.Sprintf("Nanoseconds: %d -- with value %s\n", r.When.Nanosecond(), r.Value))
		}
	}
	return fwdReadings
}

func serviceTest() {
	var readings []string
	finished := make(chan bool)
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-finished:
				return
			default:
				time.Sleep(5 * time.Second)
			case <-ticker.C:
				readings = pushReadings()
				for _, r := range readings {
					fmt.Print(r)
				}
				fmt.Println("Waiting for new results...")
			}
		}
	}()
	time.Sleep(100 * time.Second)
	ticker.Stop()
	finished <- true
	fmt.Println("Service stopped.")
}

// HTTP HANDLERS

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Here you will be shown the readings:\n")
}

// MAIN

func main() {
	s := Service{}

	s.run()

	var err error
	for i := 0; i < 3; i++ {
		err = s.createDevice(i, "Thermostat", "NULL", 1000)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	s.getDevicesList()

	serviceTest()

	time.Sleep(120 * time.Second)
	s.stop()

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}
