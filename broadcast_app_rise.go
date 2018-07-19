package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)

	// If configuration undefined, set default value
	if len(configuration.AppRiseAddress) == 0 {
		configuration.AppRiseAddress = "192.168.0.255:9999"
	}
	fmt.Printf("'%s'\n", configuration.AppRiseAddress)

	destinationAddress, _ := net.ResolveUDPAddr("udp", configuration.AppRiseAddress)
	connection, _ := net.DialUDP("udp", nil, destinationAddress)
	defer connection.Close()

	now := time.Now()
	startTime := now.UnixNano()

	var data AppCommData

	initAppMessage(&data)
	for data.AppSequenceNumber < PacketLimit {
		data.Payload = []byte("Hello")
		sendAppMessage(&data, connection)
		data.Payload = []byte("Shittydata")
		sendAppMessage(&data, connection)
	}

	now = time.Now()
	stopTime := now.UnixNano()
	fmt.Println("Datagrams sent:", data.AppSequenceNumber)
	fmt.Println("Time taken:", stopTime-startTime)
	fmt.Println("Datagrams/second:", 1000000000.0*float32(data.AppSequenceNumber)/float32(stopTime-startTime))
}
