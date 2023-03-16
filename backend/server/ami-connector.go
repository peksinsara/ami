package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/peksinsara/AMI/data"
)

func ConnectToAMI(address, username, password string, peerStatus *data.PeerStatus) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	var object data.Data

	fmt.Fprintf(conn, "Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", username, password)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		activeCalls := &data.ActiveCalls{Count: 0}
		line := scanner.Text()

		if strings.HasPrefix(line, "PeerStatus") {
			data.GetPeerStatus(line, peerStatus)
			fmt.Println("Active peers:", peerStatus.Active)
			fmt.Println("Inactive peers:", peerStatus.Inactive)
			fmt.Println()
		}

		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			if key == "Event" {
				object.Event = value
			}
			if key == "ChannelState" {
				object.ChannelState = value
			}
			if key == "Linkedid" {
				object.Linkedid = value
			}
		}
		data.HandleEvent(object, activeCalls)

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
