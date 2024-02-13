package main

import (
	"fmt"
	"net"
	"sync"
)

// perform the 3 step handshake with the master
func shakeHands(mhost string, mport int, wg *sync.WaitGroup) {
	defer wg.Done()
	ip := fmt.Sprintf("%s:%d", mhost, mport)
	conn, err := net.Dial("tcp", ip) // dont defer close this connection because it breaks step 2 of the handshake
	if err != nil {
		fmt.Println("Err connecting to master :", err)
	}
	pingMaster(conn)
	sendREPLCONF(conn)
}

func sendREPLCONF(conn net.Conn) {
	listeningPort := fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%d\r\n", port)
	_, err := conn.Write([]byte(listeningPort))
	if err != nil {
		fmt.Printf("err sending listening port : %v\n", err)
	}
	psync := "3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"
	_, err = conn.Write([]byte(psync))
	if err != nil {
		fmt.Printf("err sending psync: %v\n", err)
	}
}

func pingMaster(conn net.Conn) {
	ping := "*1\r\n$4\r\nping\r\n"
	_, err := conn.Write([]byte(ping))
	if err != nil {
		fmt.Println("Err pinging master", err)
	}
}
