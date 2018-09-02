package main

// First attempt at hub. Simple and working, but missing functionality.
import (
	"context"
	"log"
	"net"
	"syscall"

	rwf "github.com/pdxiv/gonetworktest"
)

func main() {
	startSession()
}

func startSession() {
	// Load configuration from file
	configuration := rwf.GetConfiguration(rwf.ConfigFile)

	destinationAddress, _ := net.ResolveUDPAddr("udp", configuration.HubRiseAddress)
	connection, _ := net.DialUDP("udp", nil, destinationAddress)
	defer connection.Close()

	var lc net.ListenConfig
	lc = net.ListenConfig{Control: controlOnConnSetupSoReusePort}
	// Listen to incoming UDP datagrams
	pc, err := lc.ListenPacket(context.Background(), "udp", configuration.HubSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	listenToAppAndSendHub(pc, connection)
}

func listenToAppAndSendHub(pc net.PacketConn, connection *net.UDPConn) {

	// To keep track of the expected sequence number for each app
	expectedSequenceForApp := make(map[uint64]uint64)

	var hubData rwf.HubCommData
	rwf.InitHubMessage(&hubData)
	var sinkData rwf.AppCommData
	rwf.InitAppMessage(&sinkData)
	sinkData.MasterBuffer = sinkData.MasterBuffer[0:rwf.BufferAllocationSize] // Allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(sinkData.MasterBuffer)
		// Only send a Hub message if App message is valid
		if rwf.HubDecodeAppMessage(&sinkData, &expectedSequenceForApp) {
			rwf.SendHubMessage(&sinkData, &hubData, connection)
		}
	}
}

func controlOnConnSetupSoReusePort(network string, address string, c syscall.RawConn) error {
	var operr error
	var fn = func(s uintptr) {
		operr = syscall.SetsockoptInt(int(s), syscall.SOL_SOCKET, 0xF /* syscall.SO_REUSE_PORT */, 1)
	}
	if err := c.Control(fn); err != nil {
		return err
	}
	if operr != nil {
		return operr
	}
	return nil
}
