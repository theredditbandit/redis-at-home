package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	defaultPong := "+PONG\r\n"
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	data := make([]byte, 1024) // buffer to store incoming data

    for {
		_, err = conn.Read(data)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
            break
		}
		_, err = conn.Write([]byte(defaultPong))
		if err != nil {
			fmt.Println("Error sending data:", err.Error())
		}
	}

}
