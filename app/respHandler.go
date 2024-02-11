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
		param := parsed[2]
		resp := "*2\r\n$%d\r\n%v\r\n$%d\r\n%v\r\n"
		switch param {
		case "dir":
			msg := fmt.Sprintf(resp, len(param), param, len(dirvar), dirvar)
			return []byte(msg)
		case "dbfilename":
			msg := fmt.Sprintf(resp, len(param), param, len(fName), fName)
			return []byte(msg)
		}
	case "info":
		param := parsed[1]
		var role string
		if !isSlave {
			role = "role:master"
		} else {
			role = "role:slave"
		}
		resp := "$%d\r\n%s\r\n"
		switch param {
		case "replication":
			msg := fmt.Sprintf(resp, len(role), role)
			return []byte(msg)
		}
	}
	return []byte("+OK\r\n")
}

// parseInput cleans the given resp array and returns a go array
func parseInput(d string) (string, []string) {
	semiclean := strings.Split(d, "\r\n")
	var cmd string
	var cmdFound bool
	supportedCmds := map[string]struct{}{"ping": {}, "echo": {}, "set": {}, "get": {}, "config": {}, "info": {}}
	var clean []string
	for i := 1; i < len(semiclean); i++ {
		if !strings.HasPrefix(semiclean[i], "$") {
			w := strings.ToLower(semiclean[i])
			clean = append(clean, w)
			if _, ok := supportedCmds[w]; ok && !cmdFound {
				cmd = w
				cmdFound = true
			}
		}
	}
	return cmd, clean
}
