package main

import (
	"net"
	"time"
)

func main() {
	destinationAddress, _ := net.ResolveUDPAddr("udp", "192.168.135.255:9999")
	connection, _ := net.DialUDP("udp", nil, destinationAddress)
	defer connection.Close()
	for {
		connection.Write([]byte("Hello"))
		time.Sleep(1 * time.Second)
	}
}
