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
	BroadcastAddress string
}

func main() {
	_ = startSession()
}

func startSession() error {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)

	fmt.Printf("'%s'\n", configuration.BroadcastAddress)

	// If configuration undefined, set default value
	if len(configuration.BroadcastAddress) == 0 {
		configuration.BroadcastAddress = "0.0.0.0:9999"
	}

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.BroadcastAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	fmt.Println("Listening on", configuration.BroadcastAddress)

	for {
		//simple read
		buffer := make([]byte, 512)
		pc.ReadFrom(buffer)
		fmt.Println("yeyyy, incoming udp datagram!", buffer)

		//simple write
		// pc.WriteTo([]byte("Hello from client"), addr)
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
