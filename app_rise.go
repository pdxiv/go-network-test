package main

// The purpose of this program, is to test broadcast output from App to Seq
import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)

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
		time.Sleep(100 * time.Millisecond)
	}

	now = time.Now()
	stopTime := now.UnixNano()
	fmt.Println("Datagrams sent:", data.AppSequenceNumber)
	fmt.Println("Time taken:", stopTime-startTime)
	fmt.Println("Datagrams/second:", 1000000000.0*float32(data.AppSequenceNumber)/float32(stopTime-startTime))
}
