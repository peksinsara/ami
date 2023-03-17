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
		return
	}
	defer conn.Close()

	err = wss.writeStatus(conn, wss.PeerStatus, wss.ActiveCalls)
	if err != nil {
		return
	}
	previousData := ""
	for {
		psJsonStr, err := data.PeerStatusToJSON(wss.PeerStatus)
		if err != nil {
			continue
		}
		acJsonStr, err := data.ActiveCallsToJSON(wss.ActiveCalls)
		if err != nil {
			continue
		}

		jsonStr := fmt.Sprintf(`{"status":%s, "calls":%s}`, psJsonStr, acJsonStr)
		if jsonStr != previousData {
			err = conn.WriteMessage(websocket.TextMessage, []byte(jsonStr))
			if err != nil {
				return
			}
			previousData = jsonStr
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (wss *WebSocketServer) writeStatus(conn *websocket.Conn, peerStatus *data.PeerStatus, activeCalls *data.ActiveCalls) error {
	psJsonStr, err := data.PeerStatusToJSON(peerStatus)
	if err != nil {
		return err
	}
	acJsonStr, err := data.ActiveCallsToJSON(activeCalls)
	if err != nil {
		return err
	}

	jsonStr := fmt.Sprintf(`{"status":%s, "calls":%s}`, psJsonStr, acJsonStr)
	return conn.WriteMessage(websocket.TextMessage, []byte(jsonStr))
}
