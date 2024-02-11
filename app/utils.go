package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

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
