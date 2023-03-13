package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/peksinsara/AMI/data"
)

func ConnectToAMI(address, username, password string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", username, password)

	peerStatus := &data.PeerStatus{}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if strings.HasPrefix(line, "Event: PeerStatus") {
			data.GetPeerStatus(line, peerStatus)
			fmt.Println("Updating status:", peerStatus)
			fmt.Println("Active peers:", peerStatus.Active)
			fmt.Println("Inactive peers:", peerStatus.Inactive)

		} else if strings.HasPrefix(line, "Event: CoreShowChannels") {
			activeCalls, numActiveCalls := data.GetActiveCalls(line)
			fmt.Printf("Active calls: %d\n", numActiveCalls)
			for _, call := range activeCalls {
				fmt.Println(call)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
