package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

// reads through the temperatures provided by lm-sensors
// giving every device a different reading
// on my laptop it provides 3 readings thus 3 devices
func valueService(n int) string {
	out, _ := exec.Command("sensors").Output()
	regexResult := regexp.MustCompile(".{4}.C\\s").FindAll(out, 6)
	result := string(regexResult[n])
	return result
}

func serviceTest(s *Service) {
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
				readings = s.getReadings()
				for _, r := range readings {
					fmt.Print(r)
				}
				fmt.Println("Waiting for new results...")
			}
		}
	}()
	time.Sleep(30 * time.Second)
	ticker.Stop()
	finished <- true
	s.stop()
	fmt.Println("Service stopped.")
}

// HTTP HANDLERS

/*func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Here you will be shown the readings:\n")
}*/

func main() {
	dao := &DataAccessObject{make(map[int][]DeviceReading)}
	s := Service{}

	s.init(dao)
	s.run()

	var err error
	var dev *Device
	var input *RawInput
	for i := 0; i < 3; i++ {
		input = &RawInput{
			Id:       strconv.Itoa(i),
			Name:     "Thermostat",
			Interval: "1000",
		}
		dev, err = s.createDevice(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = s.startDevice(dev, valueService)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	s.getDevicesList()

	serviceTest(&s)

	//http.HandleFunc("/", indexHandler)
	//http.ListenAndServe(":8000", nil)
}
