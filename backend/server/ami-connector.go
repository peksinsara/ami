package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/peksinsara/AMI/data"
)

func ConnectToAMI(address, username, password string, peerStatus *data.PeerStatus, activeCalls *data.ActiveCalls) error {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	var object data.Data

	fmt.Fprintf(conn, "Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", username, password)

	buf := make([]byte, 1024)
	fmt.Fprintf(conn, "Action: Command\r\n")
	fmt.Fprintf(conn, "Command: sip show peers\r\n")
	fmt.Fprintf(conn, "\r\n")
	response := ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading command response:", err)
			conn.Close()
		}
		response += string(buf[:n])
		if n < len(buf) {
			break
		}
	}
	var numOnline, numOffline int
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "OK") {
			numOnline++
		} else if strings.Contains(line, "UNKNOWN") {
			numOffline++
		}
	}

	data.InitialStatus(peerStatus, numOffline, numOnline)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		line := scanner.Text()

		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			if key == "Event" {
				object.Event = value
			} else if key == "ChannelState" {
				object.ChannelState = value
			} else if key == "Linkedid" {
				object.Linkedid = value
			} else if key == "Peer" {
				object.Peer = value
			} else if key == "PeerStatus" {
				object.PeerStatus = value
			} else {
				continue
			}

		}
		data.HandleEvent(object, activeCalls)
		data.GetPeerStatus(object, peerStatus)
		select {
		case EventNotifier <- true:
		default:
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
