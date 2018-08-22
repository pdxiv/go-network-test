package main

// The purpose of this program, is to have an App listen to Hub and respond
import (
	reuse "github.com/libp2p/go-reuseport"
	"log"
	"net"
	"time"
)

func main() {
	startSession()
}

func startSession() {
	// Load configuration from file
	configuration := GetConfiguration(ConfigFile)

	appState := InitAppState(4646)
	log.Print("Send queue has the capacity of this number of entries: ", len(appState.SendQueue))

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.AppSinkAddress)
	defer pc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize time ticker for keeping track of when events happen
	ticker := time.NewTicker(time.Nanosecond)
	defer ticker.Stop()
	latestTime := time.Now().UnixNano() // Initialize timestamp

	// Initialize channel for receiving
	appReceiver := make(chan AppCommData, 1)

	go receiveHubMessage(pc, appReceiver)
	for {
		select {
		case t := <-ticker.C:
			latestTime = t.UnixNano()
		case messageReceived := <-appReceiver:
			log.Print("Message: ", string(messageReceived.Payload), " Time: ", latestTime)
		}
	}
}

func receiveHubMessage(pc net.PacketConn, appReceiver chan AppCommData) {
	var hubData HubCommData
	InitHubMessage(&hubData)
	var appData AppCommData
	InitAppMessage(&appData)
	hubData.MasterBuffer = hubData.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer

	for {
		// Simple read
		pc.ReadFrom(hubData.MasterBuffer)
		if DecodeHubMessage(&hubData) {
			// Copy the payload of the hub message to the Master Buffer of the app message
			appData.MasterBuffer = hubData.Payload
			AppDecodeAppMessage(&appData)
			appReceiver <- appData
			hubData.ExpectedHubSequenceNumber++
		} else {
		}
	}
}
