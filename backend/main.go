package main

import (
	"fmt"

	"github.com/peksinsara/AMI/server"
)

func main() {
	fmt.Println("Connecting to AMI server...")
	server.Connect()

}
