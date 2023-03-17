package data

import (
	"encoding/json"
	"strings"
)

type PeerStatus struct {
	Active   int `json:"active_peers"`
	Inactive int `json:"inactive_peers"`
	Total    int `json:"total_peers"`
}

func InitialStatus(numOffline int, numOnline int) *PeerStatus {
	numTotal := numOffline + numOnline
	return &PeerStatus{Active: numOnline, Inactive: numOffline, Total: numTotal}

}

func (ps *PeerStatus) UpdateStatus(status string) {
	if status == "Registered" {
		if ps.Inactive > 0 {
			ps.Inactive--
		}
		ps.Active++
	} else if status == "Unregistered" {
		if ps.Active > 0 {
			ps.Active--
		}
		ps.Inactive++
	} else if status == "Reachable" {
		return
	} else {
		return
	}

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
