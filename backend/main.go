package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {

	// Connect to Asterisk Manager
	fmt.Println("Connecting to Asterisk Manager...")
	conn, err := net.Dial("tcp", "192.168.1.27:5038")
	if err != nil {
		fmt.Println("Error connecting to Asterisk Manager:", err)
		return
	}

	fmt.Println("Connected to Asterisk Manager")

	// Login to Asterisk Manager
	fmt.Fprintf(conn, "Action: Login\r\n")
	fmt.Fprintf(conn, "Username: admin\r\n")
	fmt.Fprintf(conn, "Secret: 1234\r\n")
	fmt.Fprintf(conn, "\r\n")

	// Read login response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading login response:", err)
		conn.Close()
		return
	}
	response := string(buf[:n])
	if !strings.Contains(response, "Success") {
		fmt.Println("Error logging in to Asterisk Manager:", response)
		conn.Close()

		return
	}
	fmt.Println("Logged in to Asterisk Manager")

	// Get number of peers and online/offline status
	fmt.Fprintf(conn, "Action: Command\r\n")
	fmt.Fprintf(conn, "Command: sip show peers\r\n")
	fmt.Fprintf(conn, "\r\n")

	// Read command response
	response = ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading command response:", err)
			conn.Close()

			return
		}
		response += string(buf[:n])
		if n < len(buf) {
			break
		}
	}
	var numOnline, numOffline int
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "OK (") {
			numOnline++
		} else if strings.Contains(line, "UNKNOWN") {
			numOffline++
		}
	}
	fmt.Printf("Fetched peer status at %v\n", time.Now())

	// Get channel information
	fmt.Fprintf(conn, "Action: Command\r\n")
	fmt.Fprintf(conn, "Command: core show channels\r\n")
	fmt.Fprintf(conn, "\r\n")

	// Read command response
	response = ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading command response:", err)
			conn.Close()
			return
		}
		response += string(buf[:n])
		if n < len(buf) {
			break
		}
	}
	fmt.Printf("Fetched channel data at %v\n", time.Now())

	// Count active channels, active calls, and calls processed
	var numActiveChannels, numActiveCalls, numCallsProcessed int
	words := strings.Fields(response)
	for i, word := range words {
		if word == "active" && words[i+1] == "channels" {
			fmt.Sscanf(words[i-1], "%d", &numActiveChannels)
		} else if word == "active" && words[i+1] == "call" {
			fmt.Sscanf(words[i-1], "%d", &numActiveCalls)
		} else if word == "calls" && words[i+1] == "processed" {
			fmt.Sscanf(words[i-1], "%d", &numCallsProcessed)
		}
	}
	// Print results
	fmt.Println("Total users:", numOnline+numOffline)
	fmt.Println("Online:", numOnline)
	fmt.Println("Offline:", numOffline)
	fmt.Printf("Active channels: %d\n", numActiveChannels)
	fmt.Printf("Active calls: %d\n", numActiveCalls)
	fmt.Printf("Call processed: %d\n", numCallsProcessed)

}
