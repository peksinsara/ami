package data

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	Event        string `json:"Event"`
	ChannelState string `json:"ChannelState"`
	Linkedid     string `json:"Linkedid"`
	Peer         string `json:"Peer"`
	PeerStatus   string `json:"PeerStatus"`
}

type ActiveCalls struct {
	Count int `json:"active_calls"`
}

var callIDs []string

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func removeElement(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func HandleEvent(data Data, activeCalls *ActiveCalls) {
	if data.Event == "Newstate" {
		if data.ChannelState == "6" {
			if !stringInSlice(data.Linkedid, callIDs) {
				callIDs = append(callIDs, data.Linkedid)
				activeCalls.Count++

			}
		}
	} else if data.Event == "Hangup" {
		if stringInSlice(data.Linkedid, callIDs) {
			callIDs = removeElement(callIDs, data.Linkedid)
			fmt.Println("Hangup")
			activeCalls.Count--

		}
	}
}

func ActiveCallsToJSON(activeCalls *ActiveCalls) (string, error) {
	jsonBytes, err := json.Marshal(activeCalls)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
