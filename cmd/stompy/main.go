package main

// The purpose of this program, is to have an App listen to Hub and respond
import (
	"context"
	"log"
	"net"
	"time"

	rwf "github.com/pdxiv/gonetworktest"
)

func main() {
	startSession()
}

func startSession() {
	// Load configuration from file
	configuration := rwf.GetConfiguration(rwf.ConfigFile)

	appState := rwf.InitAppState(4646)
	log.Print("Send queue has the capacity of this number of entries: ", len(appState.SendQueue))

	var lc net.ListenConfig
	lc = net.ListenConfig{Control: rwf.ControlOnConnSetupSoReusePort}
	// Listen to incoming UDP datagrams
	pc, err := lc.ListenPacket(context.Background(), "udp", configuration.AppSinkAddress)
	defer pc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize time ticker for keeping track of when events happen
	ticker := time.NewTicker(time.Nanosecond)
	defer ticker.Stop()
	latestTime := time.Now().UnixNano() // Initialize timestamp

	// Initialize channel for receiving
	appReceiver := make(chan rwf.AppCommData, 1)

	go receiveHubMessageAndDecode(pc, appReceiver)
	for {
		select {
		case t := <-ticker.C:
			latestTime = t.UnixNano()
		case messageReceived := <-appReceiver:
			log.Print("Message: ", string(messageReceived.Payload), " Time: ", latestTime)
		}
	}
}

func receiveHubMessageAndDecode(pc net.PacketConn, appReceiver chan rwf.AppCommData) {
	var hubData rwf.HubCommData
	rwf.InitHubMessage(&hubData)
	var appData rwf.AppCommData
	rwf.InitAppMessage(&appData)
	hubData.MasterBuffer = hubData.MasterBuffer[0:rwf.BufferAllocationSize] // allocate receive buffer

	for {
		// Simple read
		pc.ReadFrom(hubData.MasterBuffer)
		if rwf.DecodeHubMessage(&hubData) {
			// Copy the payload of the hub message to the Master Buffer of the app message
			appData.MasterBuffer = hubData.Payload
			rwf.AppDecodeAppMessage(&appData)
			appReceiver <- appData
			hubData.ExpectedHubSequenceNumber++
		} else {
		}
	}
}
