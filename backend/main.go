package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.27:5038")
	if err != nil {
		fmt.Println("Error connecting to Asterisk Manager:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "Action: Login\r\n")
	fmt.Fprintf(conn, "Username: admin\r\n")
	fmt.Fprintf(conn, "Secret: 1234\r\n")
	fmt.Fprintf(conn, "\r\n")

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading login response:", err)
		return
	}
	response := string(buf[:n])
	if !strings.Contains(response, "Success") {
		fmt.Println("Error logging in to Asterisk Manager:", response)
		return
	}

	fmt.Fprintf(conn, "Action: Command\r\n")
	fmt.Fprintf(conn, "Command: sip show peers\r\n")
	fmt.Fprintf(conn, "\r\n")

	response = ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading command response:", err)
			return
		}
		response += string(buf[:n])
		if n < len(buf) {
			break
		}
	}
	fmt.Println(response)

	var numOnline, numOffline int
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "OK (") {
			numOnline++
		} else if strings.Contains(line, "UNKNOWN") {
			numOffline++
		}
	}

	fmt.Println("Total users:", numOnline+numOffline)
	fmt.Println("Online:", numOnline)
	fmt.Println("Offline:", numOffline)
}
