package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func Connect() {
	conn, err := net.Dial("tcp", "192.168.1.27:5038")
	if err != nil {
		fmt.Println("Error connecting to AMI server:", err)
		return
	}

	reader := bufio.NewReader(conn)

	// Send login command
	fmt.Fprintf(conn, "Action: Login\r\nUsername: admin\r\nSecret: 1234\r\n\r\n")

	// Wait for login response
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from AMI server:", err)
			return
		}
		fmt.Println(strings.TrimSpace(line))
		if strings.HasPrefix(line, "Response: Success") {
			break
		}
	}

	// Enable event notifications
	fmt.Fprintf(conn, "Action: Events\r\nEventMask: all\r\n\r\n")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from AMI server:", err)
			return
		}
		fmt.Println(strings.TrimSpace(line))
	}
}
