package main

type Dao struct {
	Readings map[int][]DeviceReading
	Devices  map[int]Device
}
