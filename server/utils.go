package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"math/big"
	"sync"
	"time"
)

// ActiveShells is a thread safe way for us to keep track of all connected shells
var ActiveShells = sync.Map{}

// Pre-defining the ping/pong packets
var pingPacket = []byte{0x0}
var pongPacket = []byte{0x1}

const delimiter = byte('\n')

// authPayload is a weak but somewhat effective way of protecting against rouge clients
var authPayload = []byte{1, 3, 3, 7, 0x0, 4, 2, 0, 0x0, byte('l'), byte('e'), byte('e'), byte('t')}

func RandomSleepTime() time.Duration {
	// Generating a random number up to 30 and using that to sleep for the n amount of seconds
	bigIntN, _ := rand.Int(rand.Reader, big.NewInt(60))

	N := bigIntN.Int64()

	// Converting the N to time.Duration and telling it to be seconds
	return time.Duration(N) * time.Second
}

func RandomConnectionIdentifier() string {
	// Creating a 64byte buffer to randomize it later
	buf := make([]byte, 64)

	_, err := rand.Read(buf)
	if err != nil {
		log.Printf("Error: %s while generating a random identifer.\n", err.Error()); return ""
	}

	// Using SHA512 to get a random identifier, then converting it to hex
	hashFunction := sha512.New()
	hashFunction.Write(buf)
	return hex.EncodeToString(hashFunction.Sum(nil))[:16]
}
