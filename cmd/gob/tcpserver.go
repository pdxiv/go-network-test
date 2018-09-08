package main

import (
	"log"
	"net"
)

type publishMessage struct {
	message   []byte
	sessionID int
}

func startServer(listenPort string) {
	newConnections := make(chan net.Conn, 128)
	deadConnectionsIDs := make(chan int, 128)
	publishes := make(chan publishMessage, 128)
	connections := make(map[int]net.Conn)
	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		panic(err)
	}

	// New incoming connection
	go checkForNewIncomingConnections(listener, newConnections)
	connectionCounter := 0 // Used to generate session IDs
	for {
		select {
		case connection := <-newConnections:
			connections[connectionCounter] = connection
			go newConnectionSession(connection, publishes, deadConnectionsIDs, connectionCounter)
			connectionCounter++
			log.Print("Number of connections: ", len(connections))
		case deadConnectionID := <-deadConnectionsIDs:
			// Only close the connection if it hasn't been closed already
			if _, ok := connections[deadConnectionID]; ok {
				_ = connections[deadConnectionID].Close()
				delete(connections, deadConnectionID)
			}
			log.Print("Number of connections: ", len(connections))
		case publish := <-publishes:

			log.Print("Length of incoming message: ", len(publish.message))
			/* if string(publish.message) == "shit\n" {
				originatorPublish := publishMessage{message: []byte("That's a bad word!!\n"), sessionID: publish.sessionID}
				go newPublish(originatorPublish, connections[publish.sessionID], deadConnectionsIDs)
			} */
			for session, connection := range connections {
				if publish.sessionID == session {
					originatorPublish := publishMessage{message: []byte("Thanks for publishing!\n"), sessionID: publish.sessionID}
					go newPublish(originatorPublish, connection, deadConnectionsIDs)
				} else {
					go newPublish(publish, connection, deadConnectionsIDs)
				}
			}
		}
	}
	listener.Close()
}

func newPublish(publish publishMessage, connection net.Conn, deadConnectionsIDs chan int) {
	totalWritten := 0
	for totalWritten < len(publish.message) {
		writtenThisCall, err := connection.Write(publish.message[totalWritten:])
		if err != nil {
			deadConnectionsIDs <- publish.sessionID
			break
		}
		totalWritten += writtenThisCall
	}
}

func newConnectionSession(connection net.Conn, publishes chan publishMessage, deadConnectionsIDs chan int, id int) {
	buf := make([]byte, 1024)
	// Wait for incoming events
	for {
		numberOfBytes, err := connection.Read(buf)
		if err != nil {
			deadConnectionsIDs <- id
			break
		} else {
			var messageData publishMessage
			messageData.message = make([]byte, numberOfBytes)
			copy(messageData.message, buf[:numberOfBytes])
			messageData.sessionID = id
			publishes <- messageData
		}
	}
}

func checkForNewIncomingConnections(listener net.Listener, newConnections chan net.Conn) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		newConnections <- connection
	}
}
