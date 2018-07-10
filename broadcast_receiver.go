package main

import (
	"encoding/json"
	"fmt"
	reuse "github.com/libp2p/go-reuseport"
	"log"
	"os"
)

const ConfigFile = "conf.json"

type Configuration struct {
	SequencerSinkAddress string
	SequencerRiseAddress string
	AppSinkAddress       string
	AppRiseAddress       string
}

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

	for {
		// Simple read
		buffer := make([]byte, 512)
		pc.ReadFrom(buffer)
		fmt.Println("yeyyy, incoming udp datagram!", buffer)
	}
}

func getConfiguration(filename string) Configuration {
	file, _ := os.Open(filename)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
