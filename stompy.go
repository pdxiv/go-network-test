package main

// The purpose of this program, is to have an App listen to Seq and respond
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
	configuration := getConfiguration(ConfigFile)

	appState := initAppState(4646)
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

	go receiveSeqMessage(pc, appReceiver)
	for {
		select {
		case t := <-ticker.C:
			latestTime = t.UnixNano()
		case messageReceived := <-appReceiver:
			log.Print("Message: ", string(messageReceived.Payload), " Time: ", latestTime)
		}
	}
}

func receiveSeqMessage(pc net.PacketConn, appReceiver chan AppCommData) {
	var seqData SeqCommData
	initSeqMessage(&seqData)
	var appData AppCommData
	initAppMessage(&appData)
	seqData.MasterBuffer = seqData.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer

	for {
		// Simple read
		pc.ReadFrom(seqData.MasterBuffer)
		if decodeSeqMessage(&seqData) {
			// Copy the payload of the sequencer message to the Master Buffer of the app message
			appData.MasterBuffer = seqData.Payload
			appDecodeAppMessage(&appData)
			appReceiver <- appData
			seqData.ExpectedSeqSequenceNumber++
		} else {
		}
	}
}
