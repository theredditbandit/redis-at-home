package main

import (
	"fmt"
	"strings"
	"sync"
)

// respHandler responsible for reading input and generating a response
func respHandler(data []byte, dstream datastream, wg *sync.WaitGroup) []byte {
	defaultPong := "+PONG\r\n"
	d := string(data)
	cmd, parsed := parseInput(d)

	switch cmd {
	case "ping":
		return []byte(defaultPong)
	case "echo":
		msg := fmt.Sprintf("+%v\r\n", parsed[1])
		return []byte(msg)
	case "set":
		key := parsed[1]
		value := parsed[2]
		switch {
		case strings.Contains(strings.ToLower(d), "px"):
			expiry := parsed[4]
			go setExpiry(key, value, dstream, expiry, wg)
		default:
			dstream.set <- key
			dstream.set <- value
		}
		return []byte("+OK\r\n")
	case "get":
		key := parsed[1]
		dstream.get <- key
		val := <-dstream.resp
		if strings.Contains(val, "-1") {
			return []byte(val)
		}
		msg := fmt.Sprintf("+%v\r\n", val)
		return []byte(msg)
	case "config":
		fmt.Println(parsed)
	}
	return []byte("+OK\r\n")
}
