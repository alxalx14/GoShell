# GoShell

### What is GoShell?
 GoShell is a very easy and straight forward to use reverse shell.

### What is special about GoShell?
Nothing really, it supports multiple clients, has minimal stats. The client automatically daemonizes itself.

## How do I use GoShell?
To use GoShell you have 4 simple tasks to do.
1. Change the connection IP inside `client/utils.go` at line 3
2. Compile the server and the client
`go build -o server server/*.go` `go build -o client -ldflags 's -w -extdldflags "-static"' client/*.go ` (this will statically link the client to prevent a lot of compatibility issues)
3. Start the `server` script on any network that supports port forwarding.
4. Legally run the `client` script on a desired machine and that's it. The device will connect, and you can now use the CLI to execute commands on the target machine.

## Features
- Command execution network wide or single device only
- Uses `/bin/bash` as shell
- Daemonized, meaning it runs in the background
- No writes to stdout, 0 output whatsoever
- Auto reconnect until killed

# Disclaimer
GoShell is not responsible for its users actions. This software was developed out of pure boredom and was not created to be used for any kind of illegal/unethical purposes. You can perhaps use this to learn how to defend against such attacks. 

## Todo
- Support to upload/download files
- Ensure persistence using cron jobs
- Self delete after run


