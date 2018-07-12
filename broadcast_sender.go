package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

const PacketLimit = 10000 // If we're afraid of killing our network with the amount of load
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

	now := time.Now()
	startTime := now.UnixNano()
	datagramCounter := uint64(0)

	integerBuffer := make([]byte, 8)
	buf := make([]byte, 0, 65536) // Declare a byte slice send buffer with size of 64k
	for datagramCounter < PacketLimit {

		// Construct dummy "protocol" data
		buf = buf[:0] // Clear the byte slice send buffer
		binary.BigEndian.PutUint64(integerBuffer, datagramCounter)
		buf = append(buf, integerBuffer...)
		buf = append(buf, []byte("Hello")...)

		connection.Write(buf)

		datagramCounter++ // Increment every time we've sent a datagram
	}
	now = time.Now()
	stopTime := now.UnixNano()
	fmt.Println("Datagrams sent:", datagramCounter)
	fmt.Println("Time taken:", stopTime-startTime)
	fmt.Println("Datagrams/second:", 1000000000.0*float32(datagramCounter)/float32(stopTime-startTime))
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
