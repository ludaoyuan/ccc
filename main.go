package main

import (
	"log"
	"time"

	"backend"
)

func main() {
	t0 := time.Now().Unix()
	backend.Run()
	t1 := time.Now().Unix()
	log.Println(t1 - t0)
}
