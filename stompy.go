package main

// stompy.go - skeleton of a test program that just sends timestamps, to
// demonstrate what a "proper" rewheelify application should look like.

import (
	"fmt"
	"time"
)

const AppId = "STOMPY" // Wheel ID of the current application

func main() {
	// Load configuration from file
	configuration := getConfiguration(ConfigFile)
	subscribedApp := []string{"AA", "AB", "AC"}
	sender, receiver, _ := startSession(AppId, subscribedApp, configuration)

	// Initialize time ticker for keeping track of when events happen
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	latestTime := time.Now().UnixNano() // Initialize timestamp

	// Select operation to wait for incoming messages
	for {
		select {
		case messageReceived := <-receiver:
			sender <- messageReceived
		case t := <-ticker.C:
			latestTime = t.UnixNano()
			fmt.Println("Current time:", latestTime)
		}
	}
}

func startSession(selfAppId string, subscribedAppId []string, configuration Configuration) (chan<- bool, <-chan bool, bool) {
	sender := make(chan bool, 1)
	receiver := make(chan bool, 1)
	return sender, receiver, true
}
