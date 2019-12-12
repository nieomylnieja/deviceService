package main

import (
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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

func main() {
	wrtS := NewWriterService(os.Getenv("INFLUXDB_URL"),
		os.Getenv("INFLUXDB_NAME"))
	tckS := NewTickerService()
	mainS := NewService(NewDao())
	c := NewController(mainS, wrtS, tckS)

	r := newRouter(c)

	http.ListenAndServe(":8000", r)
}
