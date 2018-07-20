package main

import (
	reuse "github.com/libp2p/go-reuseport"
	"log"
)

func main() {
	_ = startSession()
}

func startSession() error {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.SequencerSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	var data AppCommData
	initAppMessage(&data)
	data.MasterBuffer = data.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(data.MasterBuffer)
		decodeAppMessage(&data)
	}
}
