package main

// The purpose of this program, is to test broadcast input from App to Hub
import (
	reuse "github.com/libp2p/go-reuseport"
	"log"
	"net"
)

func main() {
	startSession()
}

func startSession() {
	// Load configuration from file
	configuration := GetConfiguration(ConfigFile)

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.HubSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	receiveAppMessage(pc)
}

func receiveAppMessage(pc net.PacketConn) {
	var data AppCommData
	InitAppMessage(&data)
	data.MasterBuffer = data.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(data.MasterBuffer)
		HubDecodeAppMessage(&data)
	}
}
