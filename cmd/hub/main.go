package main

// First attempt at hub. Simple and working, but missing functionality.
import (
	"log"
	"net"

	reuse "github.com/libp2p/go-reuseport"
	rwf "github.com/pdxiv/gonetworktest"
)

func main() {
	startSession()
}

func startSession() {
	// Load configuration from file
	configuration := rwf.GetConfiguration(ConfigFile)

	destinationAddress, _ := net.ResolveUDPAddr("udp", configuration.HubRiseAddress)
	connection, _ := net.DialUDP("udp", nil, destinationAddress)
	defer connection.Close()

	// Listen to incoming UDP datagrams
	pc, err := reuse.ListenPacket("udp", configuration.HubSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	listenToAppAndSendHub(pc, connection)
}

func listenToAppAndSendHub(pc net.PacketConn, connection *net.UDPConn) {

	// To keep track of the expected sequence number for each app
	expectedSequenceForApp := make(map[uint64]uint64)

	var hubData HubCommData
	InitHubMessage(&hubData)
	var sinkData AppCommData
	InitAppMessage(&sinkData)
	sinkData.MasterBuffer = sinkData.MasterBuffer[0:BufferAllocationSize] // Allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(sinkData.MasterBuffer)
		// Only send a Hub message if App message is valid
		if HubDecodeAppMessage(&sinkData, &expectedSequenceForApp) {
			SendHubMessage(&sinkData, &hubData, connection)
		}
	}
}
