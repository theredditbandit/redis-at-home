package main

import (
	"fmt"
	"net"
)

// perform the 3 step handshake with the master
func shakeHands(mhost string, mport int) {
	pingMaster(mhost, mport)
}

func pingMaster(mhost string, mport int) {
	ping := "*1\r\n$4\r\nping\r\n"
	ip := fmt.Sprintf("%s:%d", mhost, mport)
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("Err connecting to master :", err)
	}
	_, err = conn.Write([]byte(ping))
	if err != nil {
		fmt.Println("Err pinging master", err)
	}
}
