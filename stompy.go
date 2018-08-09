package main

// stompy.go - skeleton of a test program that just sends timestamps, to
// demonstrate what a "proper" rewheelify application should look like.

import ()

const AppId = "STOMPY" // Wheel ID of the current application

func main() {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)
	subscribedApp := []string{"AA", "AB", "AC"}
	sender, receiver, _ := startSession(AppId, subscribedApp, configuration)

	// Select operation to wait for incoming messages
	select {
	case messageReceived := <-receiver:
		sender <- messageReceived
	}
}

func startSession(selfAppId string, subscribedAppId []string, configuration Configuration) (chan<- bool, <-chan bool, bool) {
	sender := make(chan bool, 1)
	receiver := make(chan bool, 1)
	return sender, receiver, true
}
