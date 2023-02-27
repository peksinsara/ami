package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	extensions           []string
	activeExtensions     []string
	registeredExtensions []string
	calls                int
)

func main() {
	go listenEvents()
	go serveWebsocket()
	for {
		printDashboard()
		time.Sleep(1 * time.Second)
	}
}

func listenEvents() {
	conn, err := net.Dial("tcp", "192.168.1.27:5038")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "Action: Login\r\nUsername: admin\r\nSecret: 1234\r\n\r\n")

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		message := string(buffer[:n])
		fmt.Println("Received event:", message)
		handleEvent(message)
		if strings.Contains(message, "Event: FullyBooted") {
			fmt.Println("Connected to Asterisk")
			fmt.Println("--------------------------------------")
			printDashboard()
			fmt.Println("--------------------------------------")
		}
	}
}

func handleEvent(event string) {
	if strings.Contains(event, "PeerEntry") {
		fields := strings.Split(event, "\r\n")
		peer := strings.Split(fields[1], ":")[1]
		status := strings.Split(fields[2], ":")[1]
		if strings.Contains(status, "OK") {
			extensions = appendIfMissing(extensions, peer)
			activeExtensions = appendIfMissing(activeExtensions, peer)
		} else {
			extensions = remove(extensions, peer)
			activeExtensions = remove(activeExtensions, peer)
		}
	}
	if strings.Contains(event, "PeerStatus") {
		fields := strings.Split(event, "\r\n")
		peer := strings.Split(fields[1], ":")[1]
		status := strings.Split(fields[2], ":")[1]
		if strings.Contains(status, "UNREACHABLE") {
			extensions = remove(extensions, peer)
			activeExtensions = remove(activeExtensions, peer)
		} else {
			extensions = appendIfMissing(extensions, peer)
		}
	}
	if strings.Contains(event, "Newchannel") {
		calls++
	}
	if strings.Contains(event, "Hangup") {
		calls--
	}
}

func printDashboard() {
	fmt.Printf("Total users: %d\n", len(registeredExtensions))
	fmt.Printf("Active users: %d\n", len(activeExtensions))
	fmt.Printf("Ongoing calls: %d\n", calls)
}

func serveWebsocket() {
	upgrader := websocket.Upgrader{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		for {
			// Send data to the websocket client every second
			data := fmt.Sprintf("Total users: %d\nActive users: %d\nOngoing calls: %d\n", len(extensions), len(activeExtensions), calls)
			err = conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				fmt.Println(err)
				break
			}
			time.Sleep(1 * time.Second)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func appendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func remove(slice []string, s string) []string {
	for i, ele := range slice {
		if ele == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
