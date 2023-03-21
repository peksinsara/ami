package data

import (
	"encoding/json"
	"fmt"
)

type PeerStatus struct {
	Active   int `json:"active_peers"`
	Inactive int `json:"inactive_peers"`
	Total    int `json:"total_peers"`
}

var peerAdress []string

func InitialStatus(peerStatus *PeerStatus, numOffline int, numOnline int) {
	numTotal := numOffline + numOnline
	fmt.Println("Initial status from Asterisk CLI")
	fmt.Println("Online: ", numOnline)
	fmt.Println("Offline: ", numOffline)
	fmt.Println("Total: ", numTotal)

	peerStatus.Active = numOnline
	peerStatus.Inactive = numOffline
	peerStatus.Total = numTotal
}

func GetPeerStatus(data Data, peerStatus *PeerStatus) {
	if data.Event == "PeerStatus" {
		if data.PeerStatus == "Registered" {
			if !stringInSlice(data.Peer, peerAdress) {
				peerAdress = append(peerAdress, data.Peer)
				fmt.Println("Peer registered: ", data.Peer)

				if peerStatus.Inactive > 0 {
					peerStatus.Inactive--
				}
				peerStatus.Active++

			}

		} else if data.PeerStatus == "Unregistered" {
			if stringInSlice(data.Peer, peerAdress) {
				peerAdress = removeElement(peerAdress, data.Peer)
				fmt.Println("Peer unregistered: ", data.Peer)

				if peerStatus.Active > 0 {
					peerStatus.Active--
				}

				peerStatus.Inactive++

			}
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
