package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

type datastream struct {
	set  chan string
	get  chan string
	resp chan string
}

func main() {
	ip := "0.0.0.0:6379"
	listner, err := net.Listen("tcp", ip)
	fmt.Printf("Listening on : %v\n", ip)
	var wg sync.WaitGroup
	d := datastream{
		set:  make(chan string, 2),
		get:  make(chan string),
		resp: make(chan string),
	}
	go kvHandler(d)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}
		wg.Add(1)
		go handleClient(conn, &wg, d)
	}
	wg.Wait()
}

// handleClient responsible for managing connections from multiple clients
func handleClient(conn net.Conn, wg *sync.WaitGroup, d datastream) {
	defer conn.Close()
	defer wg.Done()
	data := make([]byte, 1024) // buffer to store incoming data
	for {
		_, err := conn.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading:", err.Error())
		}
		resp := respHandler(data, d)
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println("Error sending data:", err.Error())
		}
	}
}

// respHandler responsible for reading input and generating a response
func respHandler(data []byte, dstream datastream) []byte {
	defaultPong := "+PONG\r\n"
	d := string(data)
	clean := arrayParse(d)
	if strings.Contains(strings.ToLower(d), "ping") {
		return []byte(defaultPong)
	} else if strings.Contains(strings.ToLower(d), "echo") {
		msg := fmt.Sprintf("+%v\r\n", clean[1])
		return []byte(msg)
	} else if strings.Contains(strings.ToLower(d), "set") {
		key := clean[1]
		value := clean[2]
		dstream.set <- key
		dstream.set <- value
		return []byte("+OK\r\n")
	} else if strings.Contains(strings.ToLower(d), "get") {
		key := clean[1]
		dstream.get <- key
		val := <-dstream.resp
		msg := fmt.Sprintf("+%v\r\n", val)
		return []byte(msg)
	} else {
		return []byte("+OK\r\n")
	}
}

// arrayParse cleans the given resp array and returns a go array
func arrayParse(d string) []string {
	semiclean := strings.Split(d, "\r\n")
	var clean []string
	for i := 1; i < len(semiclean); i++ {
		if !strings.HasPrefix(semiclean[i], "$") {
			clean = append(clean, semiclean[i])
		}
	}
	return clean
}

// kvHandler responsible for maintaining a dictionary of key value pairs
func kvHandler(d datastream) {
	redis := make(map[string]string)
	for {
		select {
		case k := <-d.set:
			v := <-d.set
			redis[k] = v
		case k := <-d.get:
			val, ok := redis[k]
			if ok {
				d.resp <- val
			} else {
				d.resp <- "(nil)"
			}
		}
	}

}
