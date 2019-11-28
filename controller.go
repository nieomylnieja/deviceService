package main

import "sync"

type Controller struct {
	once sync.Once
}
