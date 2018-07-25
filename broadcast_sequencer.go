package main

// Shoddy first attempt at sequencer. Some bugs and missing features right now.
import (
	"encoding/binary"
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
	receiveAppMessage(pc, connection)
}

func receiveAppMessage(pc net.PacketConn, connection *net.UDPConn) {
	var seqData SeqCommData
	initSeqMessage(&seqData)
	var sinkData AppCommData
	initAppMessage(&sinkData)
	sinkData.MasterBuffer = sinkData.MasterBuffer[0:BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(sinkData.MasterBuffer)
		decodeAppMessage(&sinkData)
		sendSeqMessage(&sinkData, &seqData, connection)
	}
}

// Encode as bytes and send a Seq message to the apps
func sendSeqMessage(sinkData *AppCommData, riseData *SeqCommData, connection *net.UDPConn) {
	// Clear riseData buffers
	riseData.MasterBuffer = riseData.MasterBuffer[:0] // Clear the byte slice send buffer

	// Convert fields into byte arrays
	binary.BigEndian.PutUint64(riseData.SessionIdBuffer, riseData.SessionId)
	binary.BigEndian.PutUint64(riseData.SeqSequenceNumberBuffer, riseData.SeqSequenceNumber)

	// Add byte arrays to master output buffer
	riseData.MasterBuffer = append(riseData.MasterBuffer, riseData.SessionIdBuffer...)
	riseData.MasterBuffer = append(riseData.MasterBuffer, riseData.SeqSequenceNumberBuffer...)

	// Add payload to master output buffer
	riseData.MasterBuffer = append(riseData.MasterBuffer, sinkData.MasterBuffer...)

	connection.Write(riseData.MasterBuffer)
	riseData.SeqSequenceNumber++ // Increment App sequence number every time we've sent a datagram
}

// Initialize all the message parameters
func initSeqMessage(data *SeqCommData) {
	data.SessionId = 0
	data.SeqSequenceNumber = 0
	data.SessionIdBuffer = make([]byte, 8)
	data.SeqSequenceNumberBuffer = make([]byte, 8)
	data.Payload = make([]byte, 0, BufferAllocationSize)
	data.MasterBuffer = make([]byte, 0, BufferAllocationSize)
}
