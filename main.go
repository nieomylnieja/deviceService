package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

// reads through the temperatures provided by lm-sensors
// giving every device a different reading
// on my laptop it provides 3 readings thus 3 devices
func valueService(n int) float64 {
	out, _ := exec.Command("sensors").Output()
	regexResult := regexp.MustCompile(".{4}.C\\s").FindAll(out, 6)
	result, _ := strconv.ParseFloat(string(regexResult[n]), 64)
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
				readings = s.GetReadings()
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
	dao := Dao{}
	s := Service{Dao: &dao}

	renv := RouterEnv{&s}
	r := renv.newRouter()

	http.ListenAndServe(":8000", r)
}
