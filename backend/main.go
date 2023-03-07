package main

import (
	"fmt"

	"github.com/peksinsara/AMI/server"
)

func main() {
	amiAddress := "192.168.7.207:5038"
	amiUsername := "admin"
	amiPassword := "1234"

	for {
		err := server.ConnectToAMI(amiAddress, amiUsername, amiPassword)
		if err != nil {
			fmt.Println("Error connecting to AMI:", err)
			return
		}
	}

}
