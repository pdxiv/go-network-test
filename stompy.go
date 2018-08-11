package main

// stompy.go - skeleton of a test program that just sends timestamps, to
// demonstrate what a "proper" rewheelify application should look like.

import (
	"log"
	"net"
	"time"

	reuse "github.com/libp2p/go-reuseport"
)

// AppID contains wheel ID of the current application
const AppID = "STOMPY"

func main() {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)
	subscribedApp := []string{"AA", "AB", "AC"}
	sender, receiver, _ := startSession(AppID, subscribedApp, configuration)

	// Initialize time ticker for keeping track of when events happen
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	latestTime := time.Now().UnixNano() // Initialize timestamp

	// Select operation to wait for incoming messages
	for {
		select {
		case t := <-ticker.C:
			latestTime = t.UnixNano()
			log.Print(latestTime)
		case messageReceived := <-receiver:
			sender <- messageReceived
		}
	}
}

func startSession(selfAppID string, subscribedAppID []string, configuration Configuration) (chan<- AppCommData, <-chan AppCommData, bool) {
	sender := make(chan AppCommData, 1)
	receiver := make(chan AppCommData, 1)

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.AppSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	receiveSeqMessage(pc)

	return sender, receiver, true
}

func receiveSeqMessage(pc net.PacketConn) {
	var data SeqCommData
	initSeqMessage(&data)
	data.MasterBuffer = data.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(data.MasterBuffer)
		if decodeSeqMessage(&data) {
			data.ExpectedSeqSequenceNumber++
			log.Print("Seq session:", data.SessionID)
			log.Print("Seq sequence:", data.SeqSequenceNumber)
			log.Print("Seq App payloads:", data.NumberOfAppPayloads)

		} else {
			log.Print("Something went wrong decoding the Seq message")
		}
	}
}
