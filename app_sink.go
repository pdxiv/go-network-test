package main

// The purpose of this program, is to test broadcast input from Seq to App
import (
	reuse "github.com/libp2p/go-reuseport"
	"log"
	"net"
)

func main() {
	startSession()
}

func startSession() {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.AppSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	receiveSeqMessage(pc)
}

func receiveSeqMessage(pc net.PacketConn) {
	var data SeqCommData
	initSeqMessage(&data)
	data.MasterBuffer = data.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(data.MasterBuffer)
		if decodeSeqMessage(&data) {
			log.Print("Got data with sequence number ", data.ExpectedSeqSequenceNumber)
			data.ExpectedSeqSequenceNumber++
		} else {
		}
	}
}
