package main

// The purpose of this program, is to test broadcast input from App to Hub
import (
	"log"
	"net"

	reuse "github.com/libp2p/go-reuseport"
	rwf "github.com/pdxiv/gonetworktest"
)

func main() {
	startSession()
}

func startSession() {
	// Load configuration from file
	configuration := rwf.GetConfiguration(rwf.ConfigFile)

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.HubSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	receiveAppMessage(pc)
}

func receiveAppMessage(pc net.PacketConn) {
	var data rwf.AppCommData
	rwf.InitAppMessage(&data)
	data.MasterBuffer = data.MasterBuffer[0:rwf.BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(data.MasterBuffer)
		rwf.HubDecodeAppMessage(&data)
	}
}
