package main

// The purpose of this program, is to test broadcast output from App to Hub
import (
	// "math/rand"
	"net"
	"time"

	rwf "github.com/pdxiv/gonetworktest"
)

func main() {
	// Load configuration from file
	configuration := rwf.GetConfiguration(rwf.ConfigFile)

	destinationAddress, _ := net.ResolveUDPAddr("udp", configuration.AppRiseAddress)
	connection, _ := net.DialUDP("udp", nil, destinationAddress)
	defer connection.Close()

	var data rwf.AppCommData

	rwf.InitAppMessage(&data)

	// Set a random dummy application ID
	//rand.Seed(time.Now().UTC().UnixNano())
	//data.ID = rand.Uint64()
	data.ID = 2323

	ticker := time.NewTicker(100 * time.Millisecond)

	for data.AppSequenceNumber < rwf.PacketLimit {
		select {
		case <-ticker.C:
			data.Payload = []byte("Hello")
			rwf.SendAppMessage(&data, connection)
		}
	}
}
