package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/peksinsara/AMI/data"
)

type WebSocketServer struct {
	PeerStatus *data.PeerStatus
}

func (wss *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}

	// Write initial status to the client
	err = wss.writeStatus(conn)
	if err != nil {
		fmt.Println("Failed to write initial status:", err)
		return
	}

	// Start a ticker to update the status every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err = wss.writeStatus(conn)
		if err != nil {
			fmt.Println("Failed to write status:", err)
			return
		}
	}
}

func (wss *WebSocketServer) writeStatus(conn *websocket.Conn) error {
	activeCalls, numActiveCalls := data.GetActiveCalls("")
	peerStatus := wss.PeerStatus

	status := struct {
		NumRegisteredPeers   int               `json:"num_registered_peers"`
		NumUnregisteredPeers int               `json:"num_unregistered_peers"`
		NumActiveCalls       int               `json:"num_active_calls"`
		ActiveCalls          []data.ActiveCall `json:"active_calls"`
	}{
		peerStatus.Active,
		peerStatus.Inactive,
		numActiveCalls,
		activeCalls,
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, jsonData)
}
