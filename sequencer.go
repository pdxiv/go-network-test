package main

// First attempt at sequencer. Simple and working, but missing functionality.
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

	destinationAddress, _ := net.ResolveUDPAddr("udp", configuration.SequencerRiseAddress)
	connection, _ := net.DialUDP("udp", nil, destinationAddress)
	defer connection.Close()

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.SequencerSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	listenToAppAndSendSeq(pc, connection)
}

func listenToAppAndSendSeq(pc net.PacketConn, connection *net.UDPConn) {

	// To keep track of the expected sequence number for each app
	expectedSequenceForApp := make(map[uint64]uint64)

	var seqData SeqCommData
	initSeqMessage(&seqData)
	var sinkData AppCommData
	initAppMessage(&sinkData)
	sinkData.MasterBuffer = sinkData.MasterBuffer[0:BufferAllocationSize] // Allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(sinkData.MasterBuffer)
		// Only send a Seq message if App message is valid
		if decodeAppMessage(&sinkData, &expectedSequenceForApp) {
			sendSeqMessage(&sinkData, &seqData, connection)
		}
	}
}
