package data

import (
	"fmt"

	"github.com/peksinsara/AMI/server"
)

func CountPeers(conn *server.AMIConn) {
	registered := 0
	unregistered := 0

	// Register event listener for PeerStatus events
	conn.RegisterEvent("PeerStatus", func(e *server.AMIEvent) {
		status := e.Fields["PeerStatus"]
		switch status {
		case "Registered":
			registered++
		case "Unregistered":
			unregistered++
		}
	})

	// Wait for events to be processed
	conn.Wait()

	// Display results
	fmt.Printf("Registered peers: %d\n", registered)
	fmt.Printf("Unregistered peers: %d\n", unregistered)
}
