package server

import (
	"fmt"
	"net"
	"strings"
)

const (
	AsteriskManagerAddress = "192.168.1.27:5038"
	AsteriskManagerUser    = "admin"
	AsteriskManagerPass    = "1234"
)

// Connect to Asterisk Manager and log in
func ConnectToAsterisk() (net.Conn, error) {
	fmt.Println("Connecting to Asterisk Manager...")
	conn, err := net.Dial("tcp", AsteriskManagerAddress)
	if err != nil {
		fmt.Println("Error connecting to Asterisk Manager:", err)
		return nil, err
	}

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
		return nil, err
	}
	response := string(buf[:n])
	if !strings.Contains(response, "Success") {
		fmt.Println("Error logging in to Asterisk Manager:", response)
		conn.Close()
		return nil, fmt.Errorf("failed to log in to Asterisk Manager")
	}

	fmt.Println("Logged in to Asterisk Manager")
	return conn, nil
}
