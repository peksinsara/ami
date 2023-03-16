package data

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	Event        string `json:"Event"`
	ChannelState string `json:"ChannelState"`
	Linkedid     string `json:"Linkedid"`
}

type ActiveCalls struct {
	Count int `json:"active_calls"`
}

func HandleEvent(data Data, activeCalls *ActiveCalls) {
	if data.Event == "Newstate" {
		fmt.Println(data.ChannelState)
		if data.ChannelState == "6" {
			activeCalls.Count++
			fmt.Println("Newstate count active calls", activeCalls.Count)
		}
	} else if data.Event == "Hangup" {
		fmt.Println(data.ChannelState)
		activeCalls.Count--
		if activeCalls.Count < 0 {
			activeCalls.Count = 0
		}
		fmt.Println("Newstate count active calls after hangup", activeCalls.Count)
	}
}

func ActiveCallsToJSON(activeCalls *ActiveCalls) (string, error) {
	jsonBytes, err := json.Marshal(activeCalls)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
