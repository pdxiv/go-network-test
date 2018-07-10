package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

const ConfigFile = "conf.json"

type Configuration struct {
	SequencerSinkAddress string
	SequencerRiseAddress string
	AppSinkAddress       string
	AppRiseAddress       string
}

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

	for {
		connection.Write([]byte("Hello"))
		time.Sleep(1 * time.Second)
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
