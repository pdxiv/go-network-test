package main

import (
	"fmt"
	reuse "github.com/libp2p/go-reuseport"
	"log"
)

func main() {
	_ = startSession()
}

func startSession() error {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)

	// If configuration undefined, set default value
	if len(configuration.SequencerSinkAddress) == 0 {
		configuration.SequencerSinkAddress = "0.0.0.0:9999"
	}
	fmt.Printf("'%s'\n", configuration.SequencerSinkAddress)

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.SequencerSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	fmt.Println("Listening on", configuration.SequencerSinkAddress)

	var data AppCommData
	initAppMessage(&data)
	data.MasterBuffer = data.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(data.MasterBuffer)
		decodeAppMessage(&data)
	}
}
