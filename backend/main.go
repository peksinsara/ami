package main

import (
	"fmt"
	"net/http"

	"github.com/peksinsara/AMI/data"
	"github.com/peksinsara/AMI/server"
)

func main() {
	amiAddress := "192.168.1.61:5038"
	amiUsername := "admin"
	amiPassword := "1234"

	// Create a new WebSocketServer instance
	wss := &server.WebSocketServer{
		PeerStatus: &data.PeerStatus{},
	}

	// Start a goroutine to serve the WebSocketServer
	go func() {
		http.Handle("/status", wss)
		err := http.ListenAndServe("192.168.1.61:8081", nil)
		if err != nil {
			fmt.Println("Error serving WebSocketServer:", err)
		}
	}()

	for {
		err := server.ConnectToAMI(amiAddress, amiUsername, amiPassword)
		if err != nil {
			fmt.Println("Error connecting to AMI:", err)
			return
		}
	}
}
