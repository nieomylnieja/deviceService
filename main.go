package main

import (
	"log"
	"net/http"
)

func main() {
	c := NewController(
		NewService(NewDao()),
		NewMeasurementsAMQP(),
		NewWriterService(),
		NewTickerService())

	r := NewRouter(c)

	log.Fatal(http.ListenAndServe(":8000", r))
}
