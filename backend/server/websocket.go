package server

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type AsteriskStatus struct {
	NumOnline         int       `json:"numOnline"`
	NumTotal          int       `json:"numTotal"`
	NumOffline        int       `json:"numOffline"`
	NumActiveChannels int       `json:"numActiveChannels"`
	NumActiveCalls    int       `json:"numActiveCalls"`
	NumCallsProcessed int       `json:"numCallsProcessed"`
	LastUpdate        time.Time `json:"lastUpdate"`
}

type Connection struct {
	conn *websocket.Conn
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request, status *AsteriskStatus) {
	// Upgrade HTTP connection to websocket connection
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		err := conn.WriteJSON(status)
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
