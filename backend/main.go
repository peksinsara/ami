package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/peksinsara/AMI/server"
)

const (
	AsteriskManagerAddress = "192.168.1.27:5038"
	AsteriskManagerUser    = "admin"
	AsteriskManagerPass    = "1234"
)

func main() {
	// Connect to Asterisk Manager
	fmt.Println("Connecting to Asterisk Manager...")
	conn, err := net.Dial("tcp", AsteriskManagerAddress)
	if err != nil {
		fmt.Println("Error connecting to Asterisk Manager:", err)
		return
	}
	defer conn.Close()

	// Login to Asterisk Manager
	fmt.Println("Connected to Asterisk Manager")
	fmt.Fprintf(conn, "Action: Login\r\n")
	fmt.Fprintf(conn, "Username: %s\r\n", AsteriskManagerUser)
	fmt.Fprintf(conn, "Secret: %s\r\n", AsteriskManagerPass)
	fmt.Fprintf(conn, "\r\n")

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

	status := &server.AsteriskStatus{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.WebsocketHandler(w, r, status)
	})
	go http.ListenAndServe(":8081", nil)

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {

		// Get total, online and offline peers
		fmt.Fprintf(conn, "Action: Command\r\n")
		fmt.Fprintf(conn, "Command: sip show peers\r\n")
		fmt.Fprintf(conn, "\r\n")

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
		// Update status
		status.NumOnline = numOnline
		status.NumOffline = numOffline
		status.NumTotal = status.NumOnline + status.NumOffline

		// Get channel info
		fmt.Fprintf(conn, "Action: Command\r\n")
		fmt.Fprintf(conn, "Command: core show channels\r\n")
		fmt.Fprintf(conn, "\r\n")
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
			} else if word == "call" && words[i+1] == "processed" {
				fmt.Sscanf(words[i-1], "%d", &numCallsProcessed)
			}
		}

		status.NumActiveChannels = numActiveChannels
		status.NumActiveCalls = numActiveCalls
		status.NumCallsProcessed = numCallsProcessed
		status.LastUpdate = time.Now()

		fmt.Printf("Total users: %d\n", status.NumOnline+status.NumOffline)
		fmt.Printf("Online: %d\n", status.NumOnline)
		fmt.Printf("Offline: %d\n", status.NumOffline)
		fmt.Printf("Active channels: %d\n", status.NumActiveChannels)
		fmt.Printf("Active calls: %d\n", status.NumActiveCalls)
		fmt.Printf("Call processed: %d\n", status.NumCallsProcessed)

	}
}
