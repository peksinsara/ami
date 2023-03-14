package data

import (
	"encoding/json"
	"strings"
)

type PeerStatus struct {
	Active   int `json:"active_peers"`
	Inactive int `json:"inactive_peers"`
}

func (ps *PeerStatus) UpdateStatus(status string) {
	switch status {
	case "Registered":
		if ps.Inactive > 0 {
			ps.Inactive--
		}
		ps.Active++
	case "Unregistered":
		if ps.Active > 0 {
			ps.Active--
		}
		ps.Inactive++
	case "Reachable":
		// do nothing, since Reachable is considered active
	default:
		// do nothing for other status
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
