package main

import (
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// Server is a structure used to create a listener for clients to join
type Server struct {
	Addr 		string
	Listener 	net.Listener
}

// Listen fires up a tcp socket server on the address specified in CLI arguments
func (s *Server) Listen() {
	listener, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatalf("Could not start up the listener: %s\n", err.Error())
	}

	s.Listener = listener

	// Waiting for clients and accepting them
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Printf("%s tried to connect. We could not accept him. Error: %s",
				conn.RemoteAddr().String(), err.Error())
			continue
		}

		// Creating a new client object
		client := new(Shell)
		// Assigning and checking the identifier
		client.Identifier = RandomConnectionIdentifier()
		if client.Identifier == "" { continue }

		client.Status = true
		client.Mutex = &sync.Mutex{}
		client.Conn = conn
		client.JoinDate = time.Now()
		// Starting the connection handler
		client.HandleConnection()
	}
}
