package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	reuse "github.com/libp2p/go-reuseport"
	"log"
	"os"
)

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

func decodeAppMessage(data *AppCommData) {
	data.PreviousAppSequenceNumber = data.AppSequenceNumber
	data.Type = binary.BigEndian.Uint16(data.MasterBuffer[0:2])
	data.PayloadSize = binary.BigEndian.Uint16(data.MasterBuffer[2:4])
	data.Id = binary.BigEndian.Uint64(data.MasterBuffer[4:12])
	data.AppSequenceNumber = binary.BigEndian.Uint64(data.MasterBuffer[12:20])
	data.Payload = data.MasterBuffer[20 : 20+data.PayloadSize]

	if data.PreviousAppSequenceNumber+1 != data.AppSequenceNumber {
		fmt.Println("**************** Gap: ", data.PreviousAppSequenceNumber+1, "to", data.AppSequenceNumber-1)
	}
	fmt.Println("Datagram type:", data.Type)
	fmt.Println("Datagram size:", data.PayloadSize)
	fmt.Println("Datagram id:", data.Id)
	fmt.Println("Datagram sequence number:", data.AppSequenceNumber)
	fmt.Printf("Datagram payload: \"%s\"\n", string(data.Payload))
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
