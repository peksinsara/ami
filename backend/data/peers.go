package data

import (
	"fmt"
	"strings"
)

type PeerStatus struct {
	Registered   int
	Unregistered int
}

func (ps *PeerStatus) UpdateStatus(status string) {
	if status == "Registered" {
		ps.Registered++
	} else if status == "Unregistered" {
		ps.Unregistered++
	}
}

func (ps *PeerStatus) String() string {
	return fmt.Sprintf("Registered peers: %d\nUnregistered peers: %d", ps.Registered, ps.Unregistered)
}

func GetPeerStatus(event string) *PeerStatus {
	peerStatus := &PeerStatus{}

	for _, line := range strings.Split(event, "\r\n") {
		if strings.HasPrefix(line, "PeerStatus: ") {
			status := strings.TrimSpace(strings.TrimPrefix(line, "PeerStatus: "))
			peerStatus.UpdateStatus(status)
		}
	}

	return peerStatus
}
