package main

// First attempt at Gob not working yet.
// Currently missing polling functionality and tcp datachannel
import (
	"context"
	"log"
	"net"

	rwf "github.com/pdxiv/gonetworktest"
)

type gobStore struct {
	data         map[uint64][][]byte
	lastSequence map[uint64]uint64
	lastSession  uint64
}

func main() {
	startSession()
}

func initGobStore(gobStorage gobStore) {
	gobStorage.data = make(map[uint64][][]byte)
	gobStorage.lastSequence = make(map[uint64]uint64)
}

func startSession() {
	// Load configuration from file
	configuration := rwf.GetConfiguration(rwf.ConfigFile)

	var gobStorage gobStore
	initGobStore(gobStorage)

	var lc net.ListenConfig
	lc = net.ListenConfig{Control: rwf.ControlOnConnSetupSoReusePort}
	// Listen to incoming UDP datagrams
	pc, err := lc.ListenPacket(context.Background(), "udp", configuration.AppSinkAddress)
	defer pc.Close()
	if err != nil {
		log.Fatal(err)
	}

	go startServer(configuration.GobTCPAddress)

	// Initialize channel for receiving
	hubReceiver := make(chan rwf.HubCommData, 1)
	messageProcessingDone := make(chan bool)
	go receiveHubMessage(pc, hubReceiver, messageProcessingDone)

	hubMasterBuffer := make([]byte, rwf.BufferAllocationSize)

	for {
		select {
		// Store raw incoming Hub messages in list, to answer Gob calls
		case messageReceived := <-hubReceiver:
			// Create session key if it doesn't exist
			if _, ok := gobStorage.data[messageReceived.SessionID]; !ok {
				gobStorage.data[messageReceived.SessionID] = make([][]byte, 0)
			}
			hubMasterBuffer = hubMasterBuffer[0:len(messageReceived.MasterBuffer)]
			copy(hubMasterBuffer, messageReceived.MasterBuffer)
			gobStorage.data[messageReceived.SessionID] = append(gobStorage.data[messageReceived.SessionID], hubMasterBuffer)

			messageProcessingDone <- true // Allow message receiver to continue, when done
		}
	}
}

func receiveHubMessage(pc net.PacketConn, hubReceiver chan rwf.HubCommData, messageProcessingDone chan bool) {
	var hubData rwf.HubCommData
	rwf.InitHubMessage(&hubData)

	hubData.MasterBuffer = hubData.MasterBuffer[0:rwf.BufferAllocationSize] // allocate receive buffer
	var frameSize int

	for {
		// Simple read
		frameSize, _, _ = pc.ReadFrom(hubData.MasterBuffer)
		hubData.MasterBuffer = hubData.MasterBuffer[0:frameSize]
		log.Print("Value of unknown integer value: ", frameSize)
		if rwf.DecodeHubMessage(&hubData) {
			log.Print("Session ID: ", hubData.SessionID)
			log.Print("Sequence number: ", hubData.HubSequenceNumber)

			hubReceiver <- hubData
			if <-messageProcessingDone {
				log.Print("Got signal that message processing is done!")
			}
			hubData.ExpectedHubSequenceNumber++
		}
	}
}
