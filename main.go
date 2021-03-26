package main

import (
	"fmt"
	"time"
)

func main() {
	timeStart := time.Now().UTC()
	for {
		time.Sleep(time.Duration(1 * time.Second))
		diffStart := time.Now().UTC().Sub(timeStart)
		fmt.Printf("time since startup: %s\n", diffStart)
	}
}
