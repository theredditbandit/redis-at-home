package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
)

type datastream struct {
	set  chan string
	get  chan string
	resp chan string
	del  chan string
}

var dirvar string
var fName string
var port int
var host = "127.0.0.1"
var mhost string // master host
var mport int    // master port
var isSlave bool

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	flag.StringVar(&dirvar, "dir", cwd, "directory")
	flag.StringVar(&fName, "dbfilename", "dump.rdb", "name of the dump file")
	flag.IntVar(&port, "port", 6379, "port")
	flag.StringVar(&mhost, "replicaof", "", "host and port which should be mirrored")
	flag.Parse()

	if len(flag.Args()) > 0 { // because flags pacakge cannot parse args like go run main.go --flag1 val0 --flag val1 val2 val3
		isSlave = true
		p, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		mport = p
	}
	if isSlave {
		shakeHands(mhost, mport)
	}
}

func main() {
	ip := fmt.Sprintf("%s:%d", host, port)
	listner, err := net.Listen("tcp", ip)
	fmt.Printf("Listening on : %v\n", ip)
	var wg sync.WaitGroup
	d := datastream{
		set:  make(chan string, 2),
		get:  make(chan string),
		resp: make(chan string),
		del:  make(chan string),
	}
	go kvHandler(d) // start the kvhandler goroutine
	if err != nil {
		fmt.Printf("Failed to bind to port %d", port)
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
		resp := respHandler(data, d, wg)
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println("Error sending data:", err.Error())
		}
	}
}

// kvHandler responsible for maintaining a dictionary of key value pairs
func kvHandler(d datastream) {
	dict := make(map[string]string)
	redisBkp := make(map[string]string)
	for {
		select {
		case k := <-d.set:
			v := <-d.set
			dict[k] = v
		case k := <-d.get:
			val, ok := dict[k]
			if ok {
				d.resp <- val
			} else {
				d.resp <- "$-1\r\n"
			}
		case k := <-d.del:
			val, ok := dict[k]
			if ok {
				redisBkp[k] = val
				delete(dict, k)
			}
		}
	}
}
