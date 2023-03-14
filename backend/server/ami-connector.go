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

	fmt.Fprintf(conn, "Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", username, password)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if strings.HasPrefix(line, "PeerStatus") {
			data.GetPeerStatus(line, peerStatus)
			fmt.Println("Active peers:", peerStatus.Active)
			fmt.Println("Inactive peers:", peerStatus.Inactive)
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
