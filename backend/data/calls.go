package data

import (
	"encoding/json"
	"strconv"
	"strings"
)

type ActiveCall struct {
	Channel            string `json:"channel"`
	CallerIDNum        string `json:"callerIDNum"`
	CallerIDName       string `json:"callerIDName"`
	Destination        string `json:"destination"`
	DestinationChannel string `json:"destinationChannel"`
}

func GetActiveCalls(event string) ([]ActiveCall, int) {
	activeCalls := []ActiveCall{}
	numActiveCalls := 0

	for _, line := range strings.Split(event, "\r\n") {
		if strings.HasPrefix(line, "Event: CoreShowChannels") {
			fields := strings.Split(line, " ")
			numActiveCalls, _ = strconv.Atoi(fields[1])
		} else if strings.HasPrefix(line, "Channel: ") {
			activeCall := ActiveCall{}
			activeCall.Channel = strings.TrimSpace(strings.TrimPrefix(line, "Channel: "))
			for _, line := range strings.Split(event, "\r\n") {
				if strings.HasPrefix(line, "ChannelState: 6") || strings.HasPrefix(line, "ChannelStateDesc: Up") {
					activeCall := ActiveCall{}
					activeCall.Channel = strings.TrimSpace(strings.TrimPrefix(line, "Channel: "))
					for _, line := range strings.Split(event, "\r\n") {
						if strings.HasPrefix(line, "CallerIDNum: ") {
							activeCall.CallerIDNum = strings.TrimSpace(strings.TrimPrefix(line, "CallerIDNum: "))
						} else if strings.HasPrefix(line, "CallerIDName: ") {
							activeCall.CallerIDName = strings.TrimSpace(strings.TrimPrefix(line, "CallerIDName: "))
						} else if strings.HasPrefix(line, "Destination: ") {
							activeCall.Destination = strings.TrimSpace(strings.TrimPrefix(line, "Destination: "))
						} else if strings.HasPrefix(line, "DestinationChannel: ") {
							activeCall.DestinationChannel = strings.TrimSpace(strings.TrimPrefix(line, "DestinationChannel: "))
						}
					}
					activeCalls = append(activeCalls, activeCall)
				}
			}
		}
	}

	return activeCalls, numActiveCalls
}

func ActiveCallsToJSON(activeCalls []ActiveCall) (string, error) {
	jsonBytes, err := json.Marshal(activeCalls)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
