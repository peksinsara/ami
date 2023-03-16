package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/peksinsara/AMI/data"
)

type WebSocketServer struct {
	PeerStatus  *data.PeerStatus
	ActiveCalls *data.ActiveCalls
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
	defer conn.Close()

	err = wss.writeStatus(conn, wss.PeerStatus)
	if err != nil {
		fmt.Println("Failed to write initial status:", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err = wss.writeStatus(conn, wss.PeerStatus)
		if err != nil {
			fmt.Println("Failed to write status:", err)
			return
		}
	}
}

func (wss *WebSocketServer) writeStatus(conn *websocket.Conn, peerStatus *data.PeerStatus) error {
	psJsonStr, err := data.PeerStatusToJSON(peerStatus)
	if err != nil {
		return err
	}

	jsonStr := fmt.Sprintf(`{"peer_status":%s,}`, psJsonStr)
	return conn.WriteMessage(websocket.TextMessage, []byte(jsonStr))
}
