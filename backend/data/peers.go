package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PeerStatus struct {
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

func (ps *PeerStatus) UpdateStatus(status string) {
	fmt.Println("Updating status:", status)
	if status == "Registered" {
		ps.Active++
	} else if status == "Unregistered" {
		ps.Inactive++
	}
}

func (ps *PeerStatus) String() string {
	return fmt.Sprintf("Active peers: %d\nInactive peers: %d\n", ps.Active, ps.Inactive)
}

func GetPeerStatus(event string, peerStatus *PeerStatus) {
	for _, line := range strings.Split(event, "\r\n") {
		if strings.HasPrefix(line, "PeerStatus: ") {
			status := strings.TrimSpace(strings.TrimPrefix(line, "PeerStatus: "))
			peerStatus.UpdateStatus(status)
		}
	}
}

func PeerStatusToJSON(peerStatus *PeerStatus) (string, error) {
	jsonBytes, err := json.Marshal(peerStatus)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
