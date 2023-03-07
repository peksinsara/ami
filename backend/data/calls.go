package data

import (
	"fmt"
	"strings"
)

type Call struct {
	Channel          string
	CallerID         string
	Destination      string
	Duration         string
	CallState        string
	AccountCode      string
	Context          string
	Exten            string
	ChannelStateDesc string
}

func GetActiveCalls(event string) ([]*Call, int) {
	var calls []*Call
	for _, line := range strings.Split(event, "\r\n") {
		if strings.HasPrefix(line, "Event: CoreShowChannels") {
			continue
		}
		fields := strings.Split(line, "!")
		if len(fields) >= 10 {
			call := &Call{
				Channel:          fields[0],
				CallerID:         fields[1],
				Destination:      fields[2],
				Duration:         fields[3],
				CallState:        fields[4],
				AccountCode:      fields[5],
				Context:          fields[6],
				Exten:            fields[7],
				ChannelStateDesc: fields[8],
			}
			calls = append(calls, call)
		}
	}

	activeCalls := 0
	for _, call := range calls {
		if call.ChannelStateDesc == "Up" {
			activeCalls++
		}
	}

	fmt.Printf("Active calls: %d\n", activeCalls)

	return calls, activeCalls
}

func (c *Call) String() string {
	return fmt.Sprintf("Channel: %s, CallerID: %s, Destination: %s, Duration: %s, CallState: %s, AccountCode: %s, Context: %s, Exten: %s, ChannelStateDesc: %s",
		c.Channel, c.CallerID, c.Destination, c.Duration, c.CallState, c.AccountCode, c.Context, c.Exten, c.ChannelStateDesc)
}
