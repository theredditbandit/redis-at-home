package main

import (
	"fmt"
	"strings"
	"sync"
	"net"
	"os"
)

func main() {
	listner, err := net.Listen("tcp", "0.0.0.0:6379")
	var wg sync.WaitGroup
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
		go handleClient(conn, &wg)
	}
	wg.Wait()
}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()
	data := make([]byte, 1024) // buffer to store incoming data
	for {
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		resp := respHandler(data)
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println("Error sending data:", err.Error())
		}
	}
}

func respHandler(data []byte) []byte {
	defaultPong := "+PONG\r\n"
    d := string(data)
	if strings.Contains(strings.ToLower(d), "ping") {
		return []byte(defaultPong)
	} else if strings.Contains(strings.ToLower(d), "echo") {
        semiclean := strings.Split(d, "\r\n")
        var clean []string
        for i := 1; i < len(semiclean); i++ {
            if !strings.HasPrefix(semiclean[i],"$") {
                clean = append(clean, semiclean[i])
            }
        }
        msg := fmt.Sprintf("+%v\r\n",clean[1])
		return []byte(msg)
	} else {
		return []byte("+OK\r\n")
	}
}
