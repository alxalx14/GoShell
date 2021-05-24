package main

import (
	"github.com/sevlyar/go-daemon"
	"log"
)

func main() {
	// Daemonizing the shell
	ctx := &daemon.Context{
		WorkDir:     "./",
	}

	proc, err := ctx.Reborn()
	if err != nil {
		log.Fatalf("Error occurred running: %s", err.Error())
	}

	if proc != nil {
		return
	}

	defer ctx.Release()

	client := new(Client)
	// Starting the connection
	client.MainLoop()
}