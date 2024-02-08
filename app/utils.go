package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// parseInput cleans the given resp array and returns a go array
func parseInput(d string) []string {
	semiclean := strings.Split(d, "\r\n")
	var clean []string
	for i := 1; i < len(semiclean); i++ {
		if !strings.HasPrefix(semiclean[i], "$") {
			clean = append(clean, semiclean[i])
		}
	}
	return clean
}

func setExpiry(key string, val string, d datastream, interval string, wg *sync.WaitGroup) {
    defer wg.Done()
	d.set <- key
	d.set <- val
	// sleep for internval miliseconds
	i, err := strconv.Atoi(interval)
	if err != nil {
		fmt.Printf("Invalid interval value: %v\n", err)
	}
	time.Sleep(time.Duration(i) * time.Millisecond)
	d.del <- key
}
