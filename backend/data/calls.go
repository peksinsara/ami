package data

import (
	"fmt"
	"strconv"
	"strings"
)

type ActiveCalls struct {
	NumCalls int `json:"calls"`
}

func GetActiveCalls(event string, activeCalls *ActiveCalls) {

	if strings.HasPrefix(event, "Event: Newstate") {
		fmt.Println("enter in HasPrefix Event: Newstate")
		activeCalls.NumCalls++
	} else if strings.HasPrefix(event, "Event: Hangup") {
		fmt.Println("enter in HasPrefix Event: Hangup")
		activeCalls.NumCalls--
	} else if strings.HasSuffix(event, "ChannelStateDesc: Up") {
		fmt.Println("enter ChannelStateDesc up ")
		activeCalls.NumCalls++
	}
}

func (ac *ActiveCalls) String() string {
	return strconv.Itoa(ac.NumCalls)
}
