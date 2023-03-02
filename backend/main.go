package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type AsteriskStatus struct {
	NumOnline         int       `json:"num_online`
	NumTotal          int       `json:"num_total`
	NumOffline        int       `json:"num_offline`
	NumActiveChannels int       `json:"num_acive_channels`
	NumActiveCalls    int       `json:"num_active_calls`
	NumCallsProcessed int       `json:"num_calls_processed`
	LastUpdate        time.Time `json:"last_update`
}

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

	// Initialize Asterisk status
	status := AsteriskStatus{}

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
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
		// Count online and offline peers
		var numOnline, numOffline int
		lines := strings.Split(response, "\n")
		for _, line := range lines {
			if strings.Contains(line, "OK (") {
				numOnline++
			} else if strings.Contains(line, "UNKNOWN") {
				numOffline++
			}
		}
		// Update status
		status.NumOnline = numOnline
		status.NumOffline = numOffline
		status.NumTotal = status.NumOnline + status.NumOffline

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
		// Update status
		status.NumActiveChannels = numActiveChannels
		status.NumActiveCalls = numActiveCalls
		status.NumCallsProcessed = numCallsProcessed

		fmt.Printf("Total users: %d\n", status.NumOnline+status.NumOffline)
		fmt.Printf("Online: %d\n", status.NumOnline)
		fmt.Printf("Offline: %d\n", status.NumOffline)
		fmt.Printf("Active channels: %d\n", status.NumActiveChannels)
		fmt.Printf("Active calls: %d\n", status.NumActiveCalls)
		fmt.Printf("Call processed: %d\n", status.NumCallsProcessed)
	}
}
