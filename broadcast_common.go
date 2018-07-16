package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
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
	ExpectedAppSequenceNumber uint64
	// Temporary buffer storage for data
	TypeBuffer              []byte
	SizeBuffer              []byte
	IdBuffer                []byte
	AppSequenceNumberBuffer []byte
	Payload                 []byte
	// The actual data as bytes that will be sent over UDP
	MasterBuffer []byte
}

// Initialize all the message parameters
func initAppMessage(data *AppCommData) {
	data.Type = 0
	data.PayloadSize = 0
	data.Id = 0
	data.AppSequenceNumber = 0
	data.ExpectedAppSequenceNumber = 0
	data.TypeBuffer = make([]byte, 2)
	data.SizeBuffer = make([]byte, 2)
	data.IdBuffer = make([]byte, 8)
	data.AppSequenceNumberBuffer = make([]byte, 8)
	data.Payload = make([]byte, 0, BufferAllocationSize)
	data.MasterBuffer = make([]byte, 0, BufferAllocationSize)
}

/*
The sequence number gap detection doesn't work properly right now.
Here's how it should work:
- At initialization, set ExpectedAppSequenceNumber to 0
- Read the App message, and decode AppSequenceNumber
- if (ExpectedAppSequenceNumber) == AppSequenceNumber then
-   increment ExpectedAppSequenceNumber
- else
-   report a sequence number gap (in the future, ask for the gaps to be filled)
-   ExpectedAppSequenceNumber = AppSequenceNumber + 1
*/
func decodeAppMessage(data *AppCommData) {
	data.ExpectedAppSequenceNumber = data.AppSequenceNumber
	data.Type = binary.BigEndian.Uint16(data.MasterBuffer[0:2])
	data.PayloadSize = binary.BigEndian.Uint16(data.MasterBuffer[2:4])
	data.Id = binary.BigEndian.Uint64(data.MasterBuffer[4:12])
	data.AppSequenceNumber = binary.BigEndian.Uint64(data.MasterBuffer[12:20])
	data.Payload = data.MasterBuffer[20 : 20+data.PayloadSize]

	if data.ExpectedAppSequenceNumber+1 != data.AppSequenceNumber {
		fmt.Println("**************** Gap: ", data.ExpectedAppSequenceNumber+1, "to", data.AppSequenceNumber-1)
	}
	fmt.Println("Datagram type:", data.Type)
	fmt.Println("Datagram size:", data.PayloadSize)
	fmt.Println("Datagram id:", data.Id)
	fmt.Println("Datagram sequence number:", data.AppSequenceNumber)
	fmt.Printf("Datagram payload: \"%s\"\n", string(data.Payload))
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
