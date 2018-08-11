package main

// The purpose of this program, is to test broadcast output from App to Seq
import (
	// "math/rand"
	"net"
	"time"
)

func main() {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)

	destinationAddress, _ := net.ResolveUDPAddr("udp", configuration.AppRiseAddress)
	connection, _ := net.DialUDP("udp", nil, destinationAddress)
	defer connection.Close()

	var data AppCommData

	initAppMessage(&data)

	// Set a random dummy application ID
	//rand.Seed(time.Now().UTC().UnixNano())
	//data.ID = rand.Uint64()
	data.ID = 2323

	ticker := time.NewTicker(100 * time.Millisecond)

	for data.AppSequenceNumber < PacketLimit {
		select {
		case <-ticker.C:
			data.Payload = []byte("Hello")
			sendAppMessage(&data, connection)
		}
	}
}
