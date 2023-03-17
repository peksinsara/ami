package main

import (
	"fmt"
	"net/http"

	"github.com/peksinsara/AMI/data"
	"github.com/peksinsara/AMI/server"
)

func main() {
	amiAddress := "192.168.1.19:5038"
	amiUsername := "admin"
	amiPassword := "1234"

	wss := &server.WebSocketServer{
		PeerStatus:  &data.PeerStatus{},
		ActiveCalls: &data.ActiveCalls{},
	}

	go func() {
		http.Handle("/status", wss)
		err := http.ListenAndServe("192.168.1.19:8082", nil)
		if err != nil {
			fmt.Println("Error serving WebSocketServer:", err)
		}
	}()

	for {
		err := server.ConnectToAMI(amiAddress, amiUsername, amiPassword, wss.PeerStatus, wss.ActiveCalls)
		if err != nil {
			fmt.Println("Error connecting to AMI:", err)
			return
		}
	}
}
