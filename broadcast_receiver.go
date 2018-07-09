package main

import (
	"fmt"
	reuse "github.com/libp2p/go-reuseport"
	"log"
	// "net"
)

func main() {
	_ = startSession()
}

func startSession() error {
	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", "192.168.0.18:9999")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		//simple read
		buffer := make([]byte, 512)
		pc.ReadFrom(buffer)
		fmt.Println("yeyyy, incoming udp datagram!", buffer)

		//simple write
		// pc.WriteTo([]byte("Hello from client"), addr)
	}
}
