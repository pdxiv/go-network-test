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
const BufferAllocationSize = 65507

type Configuration struct {
	SequencerSinkAddress string
	SequencerRiseAddress string
	AppSinkAddress       string
	AppRiseAddress       string
}

type AppCommData struct {
	// Actual data as native data types
	Type                      uint16
	PayloadSize               uint16
	Id                        uint64
	AppSequenceNumber         uint64
	PreviousAppSequenceNumber uint64 // Used to keep track of sequence number gaps
	// Temporary buffer storage for data
	TypeBuffer              []byte
	SizeBuffer              []byte
	IdBuffer                []byte
	AppSequenceNumberBuffer []byte
	Payload                 []byte
	// The actual data as bytes that will be sent over UDP
	MasterBuffer []byte
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

	var data AppCommData

	initAppMessage(&data)
	for data.AppSequenceNumber < PacketLimit {
		data.Payload = []byte("Hello")
		sendAppMessage(&data, connection)
		data.Payload = []byte("Shittydata")
		sendAppMessage(&data, connection)
	}

	now = time.Now()
	stopTime := now.UnixNano()
	fmt.Println("Datagrams sent:", data.AppSequenceNumber)
	fmt.Println("Time taken:", stopTime-startTime)
	fmt.Println("Datagrams/second:", 1000000000.0*float32(data.AppSequenceNumber)/float32(stopTime-startTime))
}

// Initialize all the message parameters
func initAppMessage(data *AppCommData) {
	data.Type = 0
	data.PayloadSize = 0
	data.Id = 0
	data.AppSequenceNumber = 0
	data.PreviousAppSequenceNumber = 0
	data.TypeBuffer = make([]byte, 2)
	data.SizeBuffer = make([]byte, 2)
	data.IdBuffer = make([]byte, 8)
	data.AppSequenceNumberBuffer = make([]byte, 8)
	data.Payload = make([]byte, 0, BufferAllocationSize)
	data.MasterBuffer = make([]byte, 0, BufferAllocationSize)
}

func sendAppMessage(data *AppCommData, connection *net.UDPConn) {
	// Clear data buffers
	data.MasterBuffer = data.MasterBuffer[:0] // Clear the byte slice send buffer

	data.PayloadSize = uint16(len(data.Payload))

	// Convert fields into byte arrays
	binary.BigEndian.PutUint16(data.TypeBuffer, data.Type)
	binary.BigEndian.PutUint16(data.SizeBuffer, data.PayloadSize)
	binary.BigEndian.PutUint64(data.IdBuffer, data.Id)
	binary.BigEndian.PutUint64(data.AppSequenceNumberBuffer, data.AppSequenceNumber)

	// Add byte arrays to master output buffer
	data.MasterBuffer = append(data.MasterBuffer, data.TypeBuffer...)
	data.MasterBuffer = append(data.MasterBuffer, data.SizeBuffer...)
	data.MasterBuffer = append(data.MasterBuffer, data.IdBuffer...)
	data.MasterBuffer = append(data.MasterBuffer, data.AppSequenceNumberBuffer...)

	// Add payload to master output buffer
	data.MasterBuffer = append(data.MasterBuffer, data.Payload...)

	connection.Write(data.MasterBuffer)
	data.AppSequenceNumber++ // Increment App sequence number every time we've sent a datagram
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
