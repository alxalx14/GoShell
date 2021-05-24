package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: ./%s <listen_address>\n", os.Args[0])
	}

	// Creating new server object
	server := new(Server)

	// Setting the listening address to the server object
	server.Addr = os.Args[1]

	// Starting the listener for the clients
	go server.Listen()

	StartUserInterface()
}
