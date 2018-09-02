package main

// The purpose of this program, is to test broadcast input from App to Hub
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

	var lc net.ListenConfig
	lc = net.ListenConfig{Control: controlOnConnSetupSoReusePort}
	// Listen to incoming UDP datagrams
	pc, err := lc.ListenPacket(context.Background(), "udp", configuration.HubSinkAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	receiveAppMessage(pc)
}

func receiveAppMessage(pc net.PacketConn) {
	var data rwf.AppCommData
	rwf.InitAppMessage(&data)
	data.MasterBuffer = data.MasterBuffer[0:rwf.BufferAllocationSize] // allocate receive buffer
	for {
		// Simple read
		pc.ReadFrom(data.MasterBuffer)
		rwf.HubDecodeAppMessage(&data)
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
