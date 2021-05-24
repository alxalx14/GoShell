package main

const server = "51.222.222.15:6969"
const delimiter = byte('\n')

// Pre-defining the ping/pong packets
var pingPacket = []byte{0x0}
var pongPacket = []byte{0x1}


// authPayload is a weak but somewhat effective way of protecting against rouge clients
var authPayload = []byte{1, 3, 3, 7, 0x0, 4, 2, 0, 0x0, byte('l'), byte('e'), byte('e'), byte('t')}



